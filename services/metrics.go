package services

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/repositories"
	"github.com/google/uuid"
)

type MetricsService interface {
	GetMonthOrdersAmount(ctx context.Context) (*models.MonthlyMetricsResponse, error)
}

type metricsService struct {
	di              *internal.Di
	orderRepository repositories.OrderRepository
}

func NewMetricsService(di *internal.Di) (MetricsService, error) {
	orderRepository, err := internal.Invoke[repositories.OrderRepository](di)
	if err != nil {
		return nil, err
	}

	return &metricsService{
		di:              di,
		orderRepository: orderRepository,
	}, nil
}

func (m *metricsService) GetMonthOrdersAmount(ctx context.Context) (*models.MonthlyMetricsResponse, error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	ordersPerMonth, err := m.orderRepository.GetOrderPerMonth(ctx, *restaurantID, nil)
	if err != nil {
		return nil, fmt.Errorf("get order per month: %w", err)
	}

	if ordersPerMonth == nil {
		return nil, models.ErrorOrderNotFound
	}

	today := time.Now().UTC()
	currentMonth := today.Format("2006-01")
	lastMonth := today.AddDate(0, -1, 0).Format("2006-01")

	var currentMonthOrders *models.OrderPerMonth
	var lastMonthOrders *models.OrderPerMonth

	for _, order := range ordersPerMonth {
		switch order.MonthWithYear {
		case currentMonth:
			currentMonthOrders = &order
		case lastMonth:
			lastMonthOrders = &order
		}
	}

	var diffFromLastMonth float64
	if currentMonthOrders != nil && lastMonthOrders != nil {
		diffFromLastMonth = (float64(currentMonthOrders.Amount) * 100 / float64(lastMonthOrders.Amount)) - 100
	}

	monthlyMetricsResponse := &models.MonthlyMetricsResponse{
		Amount:            float64(currentMonthOrders.Amount),
		DiffFromLastMonth: diffFromLastMonth,
	}

	return monthlyMetricsResponse, nil
}
