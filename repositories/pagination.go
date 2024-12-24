package repositories

import (
	"github.com/G-Villarinho/level-up-api/models"
	"gorm.io/gorm"
)

func paginate[T any](db *gorm.DB, pagination *models.Pagination, model any) (*models.PaginatedResponse[T], error) {
	var result models.PaginatedResponse[T]
	var total int64

	if err := db.Model(model).Count(&total).Error; err != nil {
		return nil, err
	}

	result.Total = total
	result.TotalPages = int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	result.Page = pagination.Page
	result.Limit = pagination.Limit

	offset := (pagination.Page - 1) * pagination.Limit
	var data []T
	query := db.Model(model).Order(pagination.Sort).Limit(pagination.Limit).Offset(offset)

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	result.Data = data
	return &result, nil
}
