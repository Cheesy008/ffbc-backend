package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogsqlc "github.com/cheesy008/ffbc-backend/internal/catalog/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

type ServiceRepository struct {
	db      *pgxpool.Pool
	queries *catalogsqlc.Queries
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		db:      db,
		queries: catalogsqlc.New(db),
	}
}

func (r *ServiceRepository) Create(
	ctx context.Context,
	input use_case.ServiceInput,
) (domain.Service, error) {
	basePrice, err := decimalToNumeric(input.BasePrice)
	if err != nil {
		return domain.Service{}, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Service{}, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	service, err := q.CreateService(ctx, catalogsqlc.CreateServiceParams{
		Name:        input.Name,
		BasePrice:   basePrice,
		Description: input.Description,
		Type:        catalogsqlc.ServiceType(input.Type),
		Status:      catalogsqlc.ServiceStatus(input.Status),
	})
	if err != nil {
		return domain.Service{}, err
	}

	for _, characteristic := range input.InputCharacteristics {
		if err := createServiceInputCharacteristic(ctx, q, int(service.ID), characteristic); err != nil {
			return domain.Service{}, err
		}
	}

	result, err := loadServiceAggregate(ctx, q, service)
	if err != nil {
		return domain.Service{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Service{}, err
	}

	return result, nil
}

func (r *ServiceRepository) Patch(
	ctx context.Context,
	input use_case.ServicePatchInput,
) (domain.Service, error) {
	params, err := mapServicePatchInput(input)
	if err != nil {
		return domain.Service{}, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Service{}, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	service, err := q.PatchService(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Service{}, domain.ErrServiceNotFound
		}
		return domain.Service{}, err
	}

	if input.Type == domain.ServiceTypeCreation && input.InputCharacteristicIDs.Set {
		if err := replaceServiceInputCharacteristics(
			ctx,
			q,
			input.ID,
			input.InputCharacteristicIDs.Value,
		); err != nil {
			return domain.Service{}, err
		}
	}

	result, err := loadServiceAggregate(ctx, q, service)
	if err != nil {
		return domain.Service{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Service{}, err
	}

	return result, nil
}

func (r *ServiceRepository) GetByID(
	ctx context.Context,
	id int,
	serviceType domain.ServiceType,
) (domain.Service, error) {
	service, err := r.queries.GetServiceByIDAndType(ctx, catalogsqlc.GetServiceByIDAndTypeParams{
		ID:   int32(id),
		Type: catalogsqlc.ServiceType(serviceType),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Service{}, domain.ErrServiceNotFound
		}
		return domain.Service{}, err
	}

	return loadServiceAggregate(ctx, r.queries, service)
}

func (r *ServiceRepository) List(
	ctx context.Context,
	filter use_case.ServiceListFilter,
) ([]domain.ServiceListItem, error) {
	services, err := r.queries.ListServices(ctx, catalogsqlc.ListServicesParams{
		Type:   catalogsqlc.ServiceType(filter.Type),
		Search: filter.Search,
		Offset: int32(filter.Offset),
		Count:  int32(filter.Count),
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.ServiceListItem, 0, len(services))
	for _, service := range services {
		item, err := mapServiceListItem(service)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *ServiceRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.queries.DeleteService(ctx, int32(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrServiceNotFound
		}
		return err
	}
	return nil
}

func replaceServiceInputCharacteristics(
	ctx context.Context,
	q *catalogsqlc.Queries,
	serviceID int,
	inputCharacteristicIDs []int,
) error {
	existing, err := q.ListServiceInputCharacteristics(ctx, int32(serviceID))
	if err != nil {
		return err
	}

	existingByID := make(map[int]use_case.ServiceInputCharacteristicInput, len(existing))
	for _, characteristic := range existing {
		existingByID[int(characteristic.ID)] = use_case.ServiceInputCharacteristicInput{
			InputCharacteristicID: int(characteristic.ID),
			IsRequired:            characteristic.IsRequired,
			SortOrder:             int32Pointer(characteristic.SortOrder),
		}
	}

	if err := q.DeleteServiceInputCharacteristics(ctx, int32(serviceID)); err != nil {
		return err
	}

	for _, inputCharacteristicID := range inputCharacteristicIDs {
		characteristic, exists := existingByID[inputCharacteristicID]
		if !exists {
			characteristic = use_case.ServiceInputCharacteristicInput{
				InputCharacteristicID: inputCharacteristicID,
				IsRequired:            true,
			}
		}
		if err := createServiceInputCharacteristic(ctx, q, serviceID, characteristic); err != nil {
			return err
		}
	}

	return nil
}

func createServiceInputCharacteristic(
	ctx context.Context,
	q *catalogsqlc.Queries,
	serviceID int,
	input use_case.ServiceInputCharacteristicInput,
) error {
	err := q.CreateServiceInputCharacteristic(ctx, catalogsqlc.CreateServiceInputCharacteristicParams{
		InputCharacteristicsID: int32(input.InputCharacteristicID),
		ServiceID:              int32(serviceID),
		IsRequired:             input.IsRequired,
		SortOrder:              nullableInt32(input.SortOrder),
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
			return domain.ErrInputCharacteristicNotFound
		}
		return err
	}
	return nil
}

func loadServiceAggregate(
	ctx context.Context,
	q *catalogsqlc.Queries,
	service catalogsqlc.Service,
) (domain.Service, error) {
	result, err := mapService(service)
	if err != nil {
		return domain.Service{}, err
	}

	if result.Type == domain.ServiceTypeCreation {
		characteristics, err := q.ListServiceInputCharacteristics(ctx, service.ID)
		if err != nil {
			return domain.Service{}, err
		}
		result.InputCharacteristics = mapServiceInputCharacteristics(characteristics)
	} else {
		result.InputCharacteristics = []domain.ServiceInputCharacteristic{}
	}

	categories, err := q.ListServiceCategoriesByServiceID(ctx, service.ID)
	if err != nil {
		return domain.Service{}, err
	}
	result.Categories = mapServiceCategories(categories)

	modifiers, err := q.ListServiceModifiersByServiceID(ctx, service.ID)
	if err != nil {
		return domain.Service{}, err
	}
	result.Modifiers = make([]domain.ServiceModifier, 0, len(modifiers))
	for _, modifier := range modifiers {
		values, err := q.ListServiceModifierValues(ctx, modifier.ID)
		if err != nil {
			return domain.Service{}, err
		}
		mappedModifier, err := mapServiceModifier(modifier, values)
		if err != nil {
			return domain.Service{}, err
		}
		result.Modifiers = append(result.Modifiers, mappedModifier)
	}

	return result, nil
}

func mapServicePatchInput(
	input use_case.ServicePatchInput,
) (catalogsqlc.PatchServiceParams, error) {
	params := catalogsqlc.PatchServiceParams{
		ID:             int32(input.ID),
		Type:           catalogsqlc.ServiceType(input.Type),
		SetName:        input.Name.Set,
		SetBasePrice:   input.BasePrice.Set,
		SetDescription: input.Description.Set,
	}

	if input.Name.Set {
		params.Name = input.Name.Value
	}
	if input.BasePrice.Set {
		basePrice, err := decimalToNumeric(input.BasePrice.Value)
		if err != nil {
			return catalogsqlc.PatchServiceParams{}, err
		}
		params.BasePrice = basePrice
	}
	if input.Description.Set {
		params.Description = input.Description.Value
	}

	return params, nil
}

func mapService(service catalogsqlc.Service) (domain.Service, error) {
	basePrice, err := numericToDecimal(service.BasePrice)
	if err != nil {
		return domain.Service{}, err
	}

	return domain.Service{
		ID:          int(service.ID),
		Name:        service.Name,
		BasePrice:   basePrice,
		Description: service.Description,
		Type:        domain.ServiceType(service.Type),
		Status:      domain.ServiceStatus(service.Status),
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
	}, nil
}

func mapServiceListItem(service catalogsqlc.ListServicesRow) (domain.ServiceListItem, error) {
	basePrice, err := numericToDecimal(service.BasePrice)
	if err != nil {
		return domain.ServiceListItem{}, err
	}

	return domain.ServiceListItem{
		ID:          int(service.ID),
		Name:        service.Name,
		BasePrice:   basePrice,
		Description: service.Description,
		Type:        domain.ServiceType(service.Type),
	}, nil
}

func mapServiceInputCharacteristics(
	rows []catalogsqlc.ListServiceInputCharacteristicsRow,
) []domain.ServiceInputCharacteristic {
	result := make([]domain.ServiceInputCharacteristic, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.ServiceInputCharacteristic{
			InputCharacteristic: domain.InputCharacteristic{
				ID:   int(row.ID),
				Name: row.Name,
				Type: domain.InputCharacteristicType(row.Type),
			},
			IsRequired: row.IsRequired,
			SortOrder:  int32Pointer(row.SortOrder),
		})
	}
	return result
}

func mapServiceCategories(rows []catalogsqlc.ServiceCategory) []domain.ServiceCategory {
	result := make([]domain.ServiceCategory, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.ServiceCategory{
			ID:        int(row.ID),
			Name:      row.Name,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}
	return result
}

func mapServiceModifier(
	modifier catalogsqlc.ServiceModifier,
	values []catalogsqlc.ServiceModifierValue,
) (domain.ServiceModifier, error) {
	mappedValues := make([]domain.ServiceModifierValue, 0, len(values))
	for _, value := range values {
		additionalPrice, err := numericToDecimal(value.AdditionalPrice)
		if err != nil {
			return domain.ServiceModifier{}, err
		}
		mappedValues = append(mappedValues, domain.ServiceModifierValue{
			ID:                int(value.ID),
			Name:              value.Name,
			ServiceModifierID: int(value.ServiceModifierID),
			AdditionalPrice:   additionalPrice,
			IsActive:          value.IsActive,
			SortOrder:         int32Pointer(value.SortOrder),
			CreatedAt:         value.CreatedAt,
			UpdatedAt:         value.UpdatedAt,
		})
	}

	return domain.ServiceModifier{
		ID:            int(modifier.ID),
		ServiceID:     int(modifier.ServiceID),
		Name:          modifier.Name,
		SelectionType: domain.ServiceModifierSelectionType(modifier.SelectionType),
		SortOrder:     int32Pointer(modifier.SortOrder),
		IsRequired:    modifier.IsRequired,
		Values:        mappedValues,
	}, nil
}

func nullableInt32(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

func int32Pointer(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	return new(int(value.Int32))
}

func decimalToNumeric(value domain.Decimal) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(value.String()); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("scan decimal %q: %w", value, err)
	}
	return numeric, nil
}

func numericToDecimal(numeric pgtype.Numeric) (domain.Decimal, error) {
	value, err := numeric.Value()
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("encode numeric: %w", err)
	}
	if value == nil {
		return decimal.Decimal{}, nil
	}

	text, ok := value.(string)
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("unexpected numeric value type %T", value)
	}
	return decimal.NewFromString(text)
}
