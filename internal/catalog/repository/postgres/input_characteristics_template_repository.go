package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogsqlc "github.com/cheesy008/ffbc-backend/internal/catalog/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

type InputCharacteristicsTemplateRepository struct {
	db      *pgxpool.Pool
	queries *catalogsqlc.Queries
}

func NewInputCharacteristicsTemplateRepository(db *pgxpool.Pool) *InputCharacteristicsTemplateRepository {
	return &InputCharacteristicsTemplateRepository{
		db:      db,
		queries: catalogsqlc.New(db),
	}
}

func (r *InputCharacteristicsTemplateRepository) Create(
	ctx context.Context,
	input use_case.InputCharacteristicsTemplateInput,
) (domain.InputCharacteristicTemplate, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	template, err := q.CreateInputCharacteristicsTemplate(ctx, catalogsqlc.CreateInputCharacteristicsTemplateParams{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	for _, inputCharacteristicID := range input.InputCharacteristicIDs {
		err := q.CreateInputCharacteristicTemplateItem(ctx, catalogsqlc.CreateInputCharacteristicTemplateItemParams{
			TemplateID:            template.ID,
			InputCharacteristicID: int32(inputCharacteristicID),
		})
		if err != nil {
			if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
				return domain.InputCharacteristicTemplate{}, domain.ErrInputCharacteristicNotFound
			}
			return domain.InputCharacteristicTemplate{}, err
		}
	}

	inputCharacteristics, err := q.ListInputCharacteristicsByTemplateID(ctx, template.ID)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	result := mapInputCharacteristicsTemplate(template)
	result.InputCharacteristics = make([]domain.InputCharacteristic, 0, len(inputCharacteristics))
	for _, characteristic := range inputCharacteristics {
		result.InputCharacteristics = append(result.InputCharacteristics, mapTemplateInputCharacteristic(characteristic))
	}

	return result, nil
}

func (r *InputCharacteristicsTemplateRepository) Patch(
	ctx context.Context,
	input use_case.InputCharacteristicsTemplatePatchInput,
) (domain.InputCharacteristicTemplate, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	params := catalogsqlc.PatchInputCharacteristicsTemplateParams{
		ID:             int32(input.ID),
		SetName:        input.Name.Set,
		SetDescription: input.Description.Set,
	}

	if input.Name.Set {
		params.Name = input.Name.Value
	}

	if input.Description.Set {
		params.Description = input.Description.Value
	}

	template, err := q.PatchInputCharacteristicsTemplate(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InputCharacteristicTemplate{}, domain.ErrInputCharacteristicsTemplateNotFound
		}

		return domain.InputCharacteristicTemplate{}, err
	}

	if input.InputCharacteristicIDs.Set {
		if err := q.DeleteInputCharacteristicTemplateItemsByTemplateID(ctx, int32(input.ID)); err != nil {
			return domain.InputCharacteristicTemplate{}, err
		}

		for _, inputCharacteristicID := range input.InputCharacteristicIDs.Value {
			err := q.CreateInputCharacteristicTemplateItem(ctx, catalogsqlc.CreateInputCharacteristicTemplateItemParams{
				TemplateID:            int32(input.ID),
				InputCharacteristicID: int32(inputCharacteristicID),
			})
			if err != nil {
				if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
					return domain.InputCharacteristicTemplate{}, domain.ErrInputCharacteristicNotFound
				}
				return domain.InputCharacteristicTemplate{}, err
			}
		}
	}

	inputCharacteristics, err := q.ListInputCharacteristicsByTemplateID(ctx, int32(input.ID))
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	result := mapInputCharacteristicsTemplate(template)
	result.InputCharacteristics = mapTemplateInputCharacteristics(inputCharacteristics)

	return result, nil
}

func (r *InputCharacteristicsTemplateRepository) GetByID(
	ctx context.Context,
	id int,
) (domain.InputCharacteristicTemplate, error) {
	template, err := r.queries.GetInputCharacteristicsTemplateByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InputCharacteristicTemplate{}, domain.ErrInputCharacteristicsTemplateNotFound
		}

		return domain.InputCharacteristicTemplate{}, err
	}

	inputCharacteristics, err := r.queries.ListInputCharacteristicsByTemplateID(ctx, int32(id))
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}

	result := mapInputCharacteristicsTemplate(template)
	result.InputCharacteristics = mapTemplateInputCharacteristics(inputCharacteristics)

	return result, nil
}

func (r *InputCharacteristicsTemplateRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.queries.DeleteInputCharacteristicsTemplate(ctx, int32(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInputCharacteristicsTemplateNotFound
		}
		return err
	}

	return nil
}

func (r *InputCharacteristicsTemplateRepository) List(
	ctx context.Context,
	filter use_case.InputCharacteristicsTemplateListFilter,
) ([]domain.InputCharacteristicTemplate, error) {
	templates, err := r.queries.ListInputCharacteristicsTemplates(ctx, catalogsqlc.ListInputCharacteristicsTemplatesParams{
		Search:                 filter.Search,
		InputCharacteristicIds: intSliceToInt32(filter.InputCharacteristicIDs),
		Offset:                 int32(filter.Offset),
		Count:                  int32(filter.Count),
		SortOrder:              filter.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.InputCharacteristicTemplate, 0, len(templates))
	for _, template := range templates {
		result = append(result, mapInputCharacteristicsTemplate(template))
	}

	return result, nil
}

func mapInputCharacteristicsTemplate(
	template catalogsqlc.InputCharacteristicsTemplate,
) domain.InputCharacteristicTemplate {
	return domain.InputCharacteristicTemplate{
		ID:          int(template.ID),
		Name:        template.Name,
		Description: template.Description,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}
}

func mapTemplateInputCharacteristic(
	characteristic catalogsqlc.InputCharacteristic,
) domain.InputCharacteristic {
	return domain.InputCharacteristic{
		ID:        int(characteristic.ID),
		Name:      characteristic.Name,
		Type:      domain.InputCharacteristicType(characteristic.Type),
		CreatedAt: characteristic.CreatedAt,
		UpdatedAt: characteristic.UpdatedAt,
	}
}

func mapTemplateInputCharacteristics(
	characteristics []catalogsqlc.InputCharacteristic,
) []domain.InputCharacteristic {
	result := make([]domain.InputCharacteristic, 0, len(characteristics))
	for _, characteristic := range characteristics {
		result = append(result, mapTemplateInputCharacteristic(characteristic))
	}
	return result
}

func intSliceToInt32(values []int) []int32 {
	result := make([]int32, 0, len(values))
	for _, value := range values {
		result = append(result, int32(value))
	}
	return result
}
