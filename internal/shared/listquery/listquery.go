package listquery

import (
	"errors"
	"strings"
)

const (
	DefaultCount = 20
	MaxCount     = 100

	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

var (
	ErrInvalidCount     = errors.New("invalid count")
	ErrInvalidOffset    = errors.New("invalid offset")
	ErrInvalidSortField = errors.New("invalid sort field")
	ErrInvalidSortOrder = errors.New("invalid sort order")
)

type Input struct {
	Search    string
	Offset    int
	Count     int
	SortBy    string
	SortOrder string
}

type Options struct {
	DefaultSortBy    string
	AllowedSortBy    []string
	DefaultSortOrder string
}

type Filter struct {
	Search    string
	Offset    int
	Count     int
	SortBy    string
	SortOrder string
}

func NormalizeNameSorted(input Input) (Filter, error) {
	return Normalize(input, Options{
		DefaultSortBy:    "name",
		AllowedSortBy:    []string{"name"},
		DefaultSortOrder: SortOrderAsc,
	})
}

func Normalize(input Input, options Options) (Filter, error) {
	offset, count, err := NormalizePagination(input.Offset, input.Count)
	if err != nil {
		return Filter{}, err
	}

	sortBy, sortOrder, err := NormalizeSorting(
		input.SortBy,
		input.SortOrder,
		options,
	)
	if err != nil {
		return Filter{}, err
	}

	return Filter{
		Search:    strings.TrimSpace(input.Search),
		Offset:    offset,
		Count:     count,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

func NormalizePagination(offset, count int) (int, int, error) {
	if offset < 0 {
		return 0, 0, ErrInvalidOffset
	}
	if count < 0 || count > MaxCount {
		return 0, 0, ErrInvalidCount
	}
	if count == 0 {
		count = DefaultCount
	}

	return offset, count, nil
}

func NormalizeSorting(sortBy, sortOrder string, options Options) (string, string, error) {
	if len(options.AllowedSortBy) == 0 {
		return "", "", nil
	}

	if sortBy == "" {
		sortBy = options.DefaultSortBy
	}
	if !contains(options.AllowedSortBy, sortBy) {
		return "", "", ErrInvalidSortField
	}

	if sortOrder == "" {
		sortOrder = options.DefaultSortOrder
	}
	if sortOrder != SortOrderAsc && sortOrder != SortOrderDesc {
		return "", "", ErrInvalidSortOrder
	}

	return sortBy, sortOrder, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
