package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"

	adminhttp "github.com/cheesy008/ffbc-backend/internal/admin/http"
	admin_postgres "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres"
	adminsqlc "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres/sqlc/generated"
	adminusecase "github.com/cheesy008/ffbc-backend/internal/admin/use_case"
	catalog_postgres "github.com/cheesy008/ffbc-backend/internal/catalog/repository/postgres"
	catalogsqlc "github.com/cheesy008/ffbc-backend/internal/catalog/repository/postgres/sqlc/generated"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/cheesy008/ffbc-backend/internal/config"
	"github.com/cheesy008/ffbc-backend/internal/storage/postgres"
)

type App struct {
	config config.Config
	server *http.Server
	db     *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := postgres.NewPool(ctx, cfg.Postgres.URL())
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	adminQueries := adminsqlc.New(db)
	adminRepo := admin_postgres.NewAdminUserRepository(adminQueries)
	sessionRepo := admin_postgres.NewAdminSessionRepository(adminQueries)
	authUseCase := adminusecase.NewAuthUseCase(adminRepo, sessionRepo)

	catalogQueries := catalogsqlc.New(db)
	categoryRepo := catalog_postgres.NewServiceCategoryRepository(catalogQueries)
	categoryUseCase := catalogusecase.NewServiceCategoryUseCase(categoryRepo)
	serviceRepo := catalog_postgres.NewServiceRepository(db)
	serviceUseCase := catalogusecase.NewServiceUseCase(serviceRepo)
	templateRepo := catalog_postgres.NewInputCharacteristicsTemplateRepository(db)
	templateUseCase := catalogusecase.NewInputCharacteristicsTemplateUseCase(templateRepo)
	inputCharacteristicsRepo := catalog_postgres.NewInputCharacteristicsRepository(db)
	inputCharacteristicsUseCase := catalogusecase.NewInputCharacteristicsUseCase(inputCharacteristicsRepo)

	mux := newRouter(authUseCase, categoryUseCase, serviceUseCase, templateUseCase, inputCharacteristicsUseCase)

	return &App{
		config: cfg,
		db:     db,
		server: &http.Server{
			Addr:              ":" + cfg.App.Port,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}, nil
}

func newRouter(
	authUseCase *adminusecase.AuthUseCase,
	categoryUseCase *catalogusecase.ServiceCategoryUseCase,
	serviceUseCase *catalogusecase.ServiceUseCase,
	templateUseCase *catalogusecase.InputCharacteristicsTemplateUseCase,
	inputCharacteristicsUseCase *catalogusecase.InputCharacteristicsUseCase,
) *http.ServeMux {
	mux := http.NewServeMux()
	api := humago.NewWithPrefix(mux, "/api", huma.DefaultConfig("FFBC Backend API", "0.1.0"))

	huma.Get(api, "/", func(ctx context.Context, input *RootInput) (*RootOutput, error) {
		return &RootOutput{
			Body: RootResponse{
				Message: "ffbc-backend api is running",
			},
		}, nil
	})
	huma.Get(api, "/healthz", func(ctx context.Context, input *HealthInput) (*HealthOutput, error) {
		return &HealthOutput{
			Body: HealthResponse{
				Status: "ok",
			},
		}, nil
	})

	adminhttp.RegisterRoutes(
		api,
		authUseCase,
		categoryUseCase,
		serviceUseCase,
		templateUseCase,
		inputCharacteristicsUseCase,
	)

	return mux
}

type RootInput struct{}

type RootOutput struct {
	Body RootResponse
}

type RootResponse struct {
	Message string `json:"message" example:"ffbc-backend api is running"`
}

type HealthInput struct{}

type HealthOutput struct {
	Body HealthResponse
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := New(ctx, cfg)
	if err != nil {
		return err
	}
	defer application.Close()

	return application.Run(ctx)
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("starting server on %s", a.server.Addr)
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}

		return nil
	}
}
