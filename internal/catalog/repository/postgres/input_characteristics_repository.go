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

type InputCharacteristicsRepository struct {
	db      *pgxpool.Pool
	queries *catalogsqlc.Queries
}

func NewInputCharacteristicsRepository(db *pgxpool.Pool) *InputCharacteristicsRepository {
	return &InputCharacteristicsRepository{
		db:      db,
		queries: catalogsqlc.New(db),
	}
}

func (r *InputCharacteristicsRepository) BulkCreate(
	ctx context.Context,
	inputs []use_case.InputCharacteristicsInput,
) ([]domain.InputCharacteristic, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	result := make([]domain.InputCharacteristic, 0, len(inputs))
	for _, input := range inputs {
		characteristic, err := q.CreateInputCharacteristic(ctx, catalogsqlc.CreateInputCharacteristicParams{
			Name: input.Name,
			Type: catalogsqlc.InputCharacteristicsType(input.Type),
		})
		if err != nil {
			return nil, err
		}

		created, err := q.GetInputCharacteristicByID(ctx, characteristic.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, mapInputCharacteristic(created))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *InputCharacteristicsRepository) Patch(
	ctx context.Context,
	input use_case.InputCharacteristicsPatchInput,
) (domain.InputCharacteristic, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.InputCharacteristic{}, err
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)
	_, err = q.PatchInputCharacteristic(ctx, mapInputCharacteristicPatchParams(input))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InputCharacteristic{}, domain.ErrInputCharacteristicNotFound
		}

		return domain.InputCharacteristic{}, err
	}

	if input.TemplateIDList.Set {
		if err := q.DeleteInputCharacteristicTemplateItems(ctx, int32(input.ID)); err != nil {
			return domain.InputCharacteristic{}, err
		}

		if input.TemplateIDList.Value != nil {
			if err := createInputCharacteristicTemplateItems(ctx, q, input.ID, *input.TemplateIDList.Value); err != nil {
				return domain.InputCharacteristic{}, mapInputCharacteristicRelationError(err)
			}
		}
	}

	result, err := q.GetInputCharacteristicByID(ctx, int32(input.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InputCharacteristic{}, domain.ErrInputCharacteristicNotFound
		}

		return domain.InputCharacteristic{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.InputCharacteristic{}, err
	}

	return mapInputCharacteristic(result), nil
}

func (r *InputCharacteristicsRepository) GetByID(ctx context.Context, id int) (domain.InputCharacteristic, error) {
	characteristic, err := r.queries.GetInputCharacteristicByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InputCharacteristic{}, domain.ErrInputCharacteristicNotFound
		}

		return domain.InputCharacteristic{}, err
	}

	return mapInputCharacteristic(characteristic), nil
}

func (r *InputCharacteristicsRepository) List(
	ctx context.Context,
	filter use_case.InputCharacteristicsListFilter,
) ([]domain.InputCharacteristic, error) {
	characteristics, err := r.queries.ListInputCharacteristics(ctx, catalogsqlc.ListInputCharacteristicsParams{
		Search: filter.Search,
		Offset: int32(filter.Offset),
		Count:  int32(filter.Count),
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.InputCharacteristic, 0, len(characteristics))
	for _, characteristic := range characteristics {
		result = append(result, mapInputCharacteristicListRow(characteristic))
	}

	return result, nil
}

func mapInputCharacteristicPatchParams(
	input use_case.InputCharacteristicsPatchInput,
) catalogsqlc.PatchInputCharacteristicParams {
	params := catalogsqlc.PatchInputCharacteristicParams{
		ID:      int32(input.ID),
		SetName: input.Name.Set,
		SetType: input.Type.Set,
		Type:    catalogsqlc.InputCharacteristicsTypeNumber,
	}

	if input.Name.Set {
		params.Name = input.Name.Value
	}

	if input.Type.Set {
		params.Type = catalogsqlc.InputCharacteristicsType(input.Type.Value)
	}

	return params
}

func createInputCharacteristicTemplateItems(
	ctx context.Context,
	q *catalogsqlc.Queries,
	inputCharacteristicID int,
	templateIDs []int,
) error {
	for _, templateID := range templateIDs {
		if err := q.CreateInputCharacteristicTemplateItem(ctx, catalogsqlc.CreateInputCharacteristicTemplateItemParams{
			TemplateID:            int32(templateID),
			InputCharacteristicID: int32(inputCharacteristicID),
		}); err != nil {
			return err
		}
	}

	return nil
}

func mapInputCharacteristicRelationError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
		return domain.ErrInputCharacteristicsTemplateNotFound
	}

	return err
}

func mapInputCharacteristic(row catalogsqlc.GetInputCharacteristicByIDRow) domain.InputCharacteristic {
	return domain.InputCharacteristic{
		ID:          int(row.ID),
		Name:        row.Name,
		Type:        domain.InputCharacteristicType(row.Type),
		TemplateIDs: int32SliceToInt(row.TemplateIds),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapInputCharacteristicListRow(row catalogsqlc.InputCharacteristic) domain.InputCharacteristic {
	return domain.InputCharacteristic{
		ID:        int(row.ID),
		Name:      row.Name,
		Type:      domain.InputCharacteristicType(row.Type),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func int32SliceToInt(values []int32) []int {
	result := make([]int, 0, len(values))
	for _, value := range values {
		result = append(result, int(value))
	}
	return result
}
