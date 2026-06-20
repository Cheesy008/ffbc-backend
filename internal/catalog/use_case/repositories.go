package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

type ServiceRepository interface {
	Create(ctx context.Context, input ServiceInput) (domain.Service, error)
	Patch(ctx context.Context, input ServicePatchInput) (domain.Service, error)
	GetByID(ctx context.Context, id int, serviceType domain.ServiceType) (domain.Service, error)
	List(ctx context.Context, filter ServiceListFilter) ([]domain.ServiceListItem, error)
	Delete(ctx context.Context, id int) error
}

type InputCharacteristicsTemplateRepository interface {
	Create(ctx context.Context, input InputCharacteristicsTemplateInput) (domain.InputCharacteristicTemplate, error)
	Patch(ctx context.Context, input InputCharacteristicsTemplatePatchInput) (domain.InputCharacteristicTemplate, error)
	GetByID(ctx context.Context, id int) (domain.InputCharacteristicTemplate, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter InputCharacteristicsTemplateListFilter) ([]domain.InputCharacteristicTemplate, error)
}

type InputCharacteristicsRepository interface {
	BulkCreate(ctx context.Context, inputs []InputCharacteristicsInput) ([]domain.InputCharacteristic, error)
	Patch(ctx context.Context, input InputCharacteristicsPatchInput) (domain.InputCharacteristic, error)
	GetByID(ctx context.Context, id int) (domain.InputCharacteristic, error)
	List(ctx context.Context, filter InputCharacteristicsListFilter) ([]domain.InputCharacteristic, error)
}
