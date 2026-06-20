package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	servicehandler "github.com/cheesy008/ffbc-backend/internal/admin/http/handlers/catalog/service"
)

func registerServiceRoutes(group *huma.Group, handler *servicehandler.Handler) {
	huma.Register(group, huma.Operation{
		OperationID: "admin-list-creation-services",
		Method:      http.MethodGet,
		Path:        "/catalog/service/creation",
		Summary:     "List creation services",
		Tags:        []string{"services"},
	}, handler.ListCreation)

	huma.Register(group, huma.Operation{
		OperationID: "admin-create-creation-service",
		Method:      http.MethodPost,
		Path:        "/catalog/service/creation",
		Summary:     "Create creation service",
		Tags:        []string{"services"},
	}, handler.CreateCreation)

	huma.Register(group, huma.Operation{
		OperationID: "admin-creation-service-detail",
		Method:      http.MethodGet,
		Path:        "/catalog/service/creation/{id}",
		Summary:     "Creation service detail",
		Tags:        []string{"services"},
	}, handler.DetailCreation)

	huma.Register(group, huma.Operation{
		OperationID: "admin-update-creation-service",
		Method:      http.MethodPatch,
		Path:        "/catalog/service/creation/{id}",
		Summary:     "Update creation service",
		Tags:        []string{"services"},
	}, handler.PatchCreation)

	huma.Register(group, huma.Operation{
		OperationID: "admin-list-selling-services",
		Method:      http.MethodGet,
		Path:        "/catalog/service/selling",
		Summary:     "List selling services",
		Tags:        []string{"services"},
	}, handler.ListSelling)

	huma.Register(group, huma.Operation{
		OperationID: "admin-create-selling-service",
		Method:      http.MethodPost,
		Path:        "/catalog/service/selling",
		Summary:     "Create selling service",
		Tags:        []string{"services"},
	}, handler.CreateSelling)

	huma.Register(group, huma.Operation{
		OperationID: "admin-selling-service-detail",
		Method:      http.MethodGet,
		Path:        "/catalog/service/selling/{id}",
		Summary:     "Selling service detail",
		Tags:        []string{"services"},
	}, handler.DetailSelling)

	huma.Register(group, huma.Operation{
		OperationID: "admin-update-selling-service",
		Method:      http.MethodPatch,
		Path:        "/catalog/service/selling/{id}",
		Summary:     "Update selling service",
		Tags:        []string{"services"},
	}, handler.PatchSelling)

	huma.Register(group, huma.Operation{
		OperationID:   "admin-delete-service",
		Method:        http.MethodDelete,
		Path:          "/catalog/service/{id}",
		Summary:       "Delete service",
		Tags:          []string{"services"},
		DefaultStatus: http.StatusNoContent,
	}, handler.Delete)
}
