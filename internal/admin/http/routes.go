package http

import (
	"net/http"

	"github.com/cheesy008/ffbc-backend/internal/admin/http/handlers/catalog"
	"github.com/cheesy008/ffbc-backend/internal/admin/http/handlers/catalog/input_characteristics"
	servicehandler "github.com/cheesy008/ffbc-backend/internal/admin/http/handlers/catalog/service"
	"github.com/danielgtaylor/huma/v2"

	"github.com/cheesy008/ffbc-backend/internal/admin/http/handlers"
	adminusecase "github.com/cheesy008/ffbc-backend/internal/admin/use_case"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

const Prefix = "/admin"

func RegisterRoutes(
	api huma.API,
	authUseCase *adminusecase.AuthUseCase,
	categoryUseCase *catalogusecase.ServiceCategoryUseCase,
	serviceUseCase *catalogusecase.ServiceUseCase,
	templateUseCase *catalogusecase.InputCharacteristicsTemplateUseCase,
	inputCharacteristicsUseCase *catalogusecase.InputCharacteristicsUseCase,
) {
	publicGroup := huma.NewGroup(api, Prefix)
	protectedGroup := huma.NewGroup(api, Prefix)
	protectedGroup.UseMiddleware(sessionMiddleware(authUseCase))

	authHandler := handlers.NewAuthHandler(authUseCase)
	adminHandler := handlers.NewAdminHandler()
	categoryHandler := catalog.NewServiceCategoryHandler(categoryUseCase)
	serviceHandler := servicehandler.New(serviceUseCase)
	templateHandler := catalog.NewInputCharacteristicsTemplateHandler(templateUseCase)
	inputCharacteristicsHandler := inputcharacteristics.New(inputCharacteristicsUseCase)

	huma.Register(publicGroup, huma.Operation{
		OperationID: "admin-login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Admin login",
		Tags:        []string{"admin-auth"},
	}, authHandler.Login)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-logout",
		Method:      http.MethodPost,
		Path:        "/auth/logout",
		Summary:     "Admin logout",
		Tags:        []string{"admin-auth"},
	}, authHandler.Logout)
	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-me",
		Method:      http.MethodGet,
		Path:        "/user/me",
		Summary:     "Current admin",
		Tags:        []string{"admin-user"},
	}, adminHandler.Me)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-create-service-category",
		Method:      http.MethodPost,
		Path:        "/catalog/service-categories",
		Summary:     "Create service category",
		Tags:        []string{"service-categories"},
	}, categoryHandler.Create)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-update-service-category",
		Method:      http.MethodPatch,
		Path:        "/catalog/service-categories/{id}",
		Summary:     "Update service category",
		Tags:        []string{"service-categories"},
	}, categoryHandler.Update)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-list-service-categories",
		Method:      http.MethodGet,
		Path:        "/catalog/service-categories",
		Summary:     "List service categories",
		Tags:        []string{"service-categories"},
	}, categoryHandler.List)

	huma.Register(protectedGroup, huma.Operation{
		OperationID:   "admin-delete-service-category",
		Method:        http.MethodDelete,
		Path:          "/catalog/service-categories/{id}",
		Summary:       "Delete service category",
		Tags:          []string{"service-categories"},
		DefaultStatus: http.StatusNoContent,
	}, categoryHandler.Delete)

	registerServiceRoutes(protectedGroup, serviceHandler)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-create-input-characteristics-template",
		Method:      http.MethodPost,
		Path:        "/catalog/input-characteristics-templates",
		Summary:     "Create input characteristics template",
		Tags:        []string{"input-characteristics-templates"},
	}, templateHandler.Create)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-update-input-characteristics-template",
		Method:      http.MethodPatch,
		Path:        "/catalog/input-characteristics-templates/{id}",
		Summary:     "Update input characteristics template",
		Tags:        []string{"input-characteristics-templates"},
	}, templateHandler.Patch)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-input-characteristics-template-detail",
		Method:      http.MethodGet,
		Path:        "/catalog/input-characteristics-templates/{id}",
		Summary:     "Input characteristics template detail",
		Tags:        []string{"input-characteristics-templates"},
	}, templateHandler.Detail)

	huma.Register(protectedGroup, huma.Operation{
		OperationID:   "admin-delete-input-characteristics-template",
		Method:        http.MethodDelete,
		Path:          "/catalog/input-characteristics-templates/{id}",
		Summary:       "Delete input characteristics template",
		Tags:          []string{"input-characteristics-templates"},
		DefaultStatus: http.StatusNoContent,
	}, templateHandler.Delete)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-list-input-characteristics-templates",
		Method:      http.MethodGet,
		Path:        "/catalog/input-characteristics-templates",
		Summary:     "List input characteristics templates",
		Tags:        []string{"input-characteristics-templates"},
	}, templateHandler.List)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-bulk-create-input-characteristic",
		Method:      http.MethodPost,
		Path:        "/catalog/input-characteristics/bulk",
		Summary:     "Bulk create input characteristics",
		Tags:        []string{"input-characteristics"},
	}, inputCharacteristicsHandler.BulkCreate)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-update-input-characteristic",
		Method:      http.MethodPatch,
		Path:        "/catalog/input-characteristics/{id}",
		Summary:     "Update input characteristic",
		Tags:        []string{"input-characteristics"},
	}, inputCharacteristicsHandler.Patch)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-input-characteristic-detail",
		Method:      http.MethodGet,
		Path:        "/catalog/input-characteristics/{id}",
		Summary:     "Input characteristic detail",
		Tags:        []string{"input-characteristics"},
	}, inputCharacteristicsHandler.Detail)

	huma.Register(protectedGroup, huma.Operation{
		OperationID: "admin-list-input-characteristics",
		Method:      http.MethodGet,
		Path:        "/catalog/input-characteristics",
		Summary:     "List input characteristics",
		Tags:        []string{"input-characteristics"},
	}, inputCharacteristicsHandler.List)
}
