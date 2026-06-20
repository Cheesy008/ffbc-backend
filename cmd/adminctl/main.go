package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	admin_postgres "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres"
	adminsqlc "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/admin/use_case"
	"github.com/cheesy008/ffbc-backend/internal/config"
	"github.com/cheesy008/ffbc-backend/internal/security/password"
	"github.com/cheesy008/ffbc-backend/internal/storage/postgres"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: adminctl create-admin --email x --password x")
	}

	switch args[0] {
	case "create-admin":
		return createAdmin(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func createAdmin(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create-admin", flag.ExitOnError)

	email := fs.String("email", "", "admin email")
	plainPassword := fs.String("password", "", "admin password")
	displayName := fs.String("display-name", "", "admin display name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *email == "" || *plainPassword == "" {
		return fmt.Errorf("email and password are required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := postgres.NewPool(ctx, cfg.Postgres.URL())
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := adminsqlc.New(pool)
	adminRepo := admin_postgres.NewAdminUserRepository(queries)
	adminUseCase := use_case.NewAdminUseCase(adminRepo)

	hashedPassword, err := password.Hash(*plainPassword)
	if err != nil {
		return err
	}

	var displayNamePtr *string
	if *displayName != "" {
		displayNamePtr = displayName
	}

	admin, err := adminUseCase.Create(ctx, use_case.AdminUserCreateInput{
		Email:       *email,
		Password:    hashedPassword,
		DisplayName: displayNamePtr,
		IsActive:    true,
	})
	if err != nil {
		return err
	}

	fmt.Printf("created admin user: id=%d email=%s\n", admin.ID, admin.Email)
	return nil
}
