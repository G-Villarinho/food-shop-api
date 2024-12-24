package models

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidPageParameter  = errors.New("invalid page parameter")
	ErrInvalidLimitParameter = errors.New("invalid limit parameter")
)

type Pagination struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func NewPagination(pageStr, limitStr, sort string) (*Pagination, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, ErrInvalidPageParameter
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return nil, ErrInvalidLimitParameter
	}

	if sort == "" {
		sort = "id ASC"
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}, nil
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

func MapPaginatedResult[T any, U any](result *PaginatedResponse[T], mapper func(T) U) *PaginatedResponse[U] {
	newData := make([]U, len(result.Data))
	for i, item := range result.Data {
		newData[i] = mapper(item)
	}

	return &PaginatedResponse[U]{
		Data:       newData,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		Page:       result.Page,
		Limit:      result.Limit,
	}
}
