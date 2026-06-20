package postgres

import (
	"context"
	"errors"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogsqlc "github.com/cheesy008/ffbc-backend/internal/catalog/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ServiceCategoryRepository struct {
	queries *catalogsqlc.Queries
}

func NewServiceCategoryRepository(queries *catalogsqlc.Queries) *ServiceCategoryRepository {
	return &ServiceCategoryRepository{queries: queries}
}

func (r *ServiceCategoryRepository) Create(
	ctx context.Context,
	input use_case.ServiceCategoryInput,
) (domain.ServiceCategory, error) {
	category, err := r.queries.CreateServiceCategory(ctx, input.Name)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ServiceCategory{}, domain.ErrServiceCategoryAlreadyExists
		}

		return domain.ServiceCategory{}, err
	}

	return mapServiceCategory(category), nil
}

func (r *ServiceCategoryRepository) Update(
	ctx context.Context,
	input use_case.ServiceCategoryInput,
) (domain.ServiceCategory, error) {
	category, err := r.queries.UpdateServiceCategory(ctx, catalogsqlc.UpdateServiceCategoryParams{
		ID:   int32(input.ID),
		Name: input.Name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ServiceCategory{}, domain.ErrServiceCategoryNotFound
		}
		if isUniqueViolation(err) {
			return domain.ServiceCategory{}, domain.ErrServiceCategoryAlreadyExists
		}

		return domain.ServiceCategory{}, err
	}

	return mapServiceCategory(category), nil
}

func (r *ServiceCategoryRepository) List(
	ctx context.Context,
	filter use_case.ServiceCategoryListFilter,
) ([]domain.ServiceCategory, error) {
	categories, err := r.queries.ListServiceCategories(ctx, catalogsqlc.ListServiceCategoriesParams{
		Search:    filter.Search,
		SortOrder: filter.SortOrder,
		Offset:    int32(filter.Offset),
		Count:     int32(filter.Count),
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.ServiceCategory, 0, len(categories))
	for _, category := range categories {
		result = append(result, mapServiceCategory(category))
	}

	return result, nil
}

func (r *ServiceCategoryRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.queries.DeleteServiceCategory(ctx, int32(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrServiceCategoryNotFound
		}
		return err
	}

	return nil
}

func mapServiceCategory(category catalogsqlc.ServiceCategory) domain.ServiceCategory {
	return domain.ServiceCategory{
		ID:        int(category.ID),
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
