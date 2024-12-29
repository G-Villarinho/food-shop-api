// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/G-Villarinho/food-shop-api/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// CreateOrderWithItems provides a mock function with given fields: ctx, order, items
func (_m *OrderRepository) CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	ret := _m.Called(ctx, order, items)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrderWithItems")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Order, []models.OrderItem) error); ok {
		r0 = rf(ctx, order, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOrderByID provides a mock function with given fields: ctx, orderID, preload
func (_m *OrderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID, preload bool) (*models.Order, error) {
	ret := _m.Called(ctx, orderID, preload)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderByID")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) (*models.Order, error)); ok {
		return rf(ctx, orderID, preload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) *models.Order); ok {
		r0 = rf(ctx, orderID, preload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, bool) error); ok {
		r1 = rf(ctx, orderID, preload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaginatedOrdersByRestaurantID provides a mock function with given fields: ctx, restaurantID, pagination
func (_m *OrderRepository) GetPaginatedOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.OrderPagination) (*models.PaginatedResponse[models.Order], error) {
	ret := _m.Called(ctx, restaurantID, pagination)

	if len(ret) == 0 {
		panic("no return value specified for GetPaginatedOrdersByRestaurantID")
	}

	var r0 *models.PaginatedResponse[models.Order]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *models.OrderPagination) (*models.PaginatedResponse[models.Order], error)); ok {
		return rf(ctx, restaurantID, pagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *models.OrderPagination) *models.PaginatedResponse[models.Order]); ok {
		r0 = rf(ctx, restaurantID, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PaginatedResponse[models.Order])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, *models.OrderPagination) error); ok {
		r1 = rf(ctx, restaurantID, pagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, orderID, status
func (_m *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	ret := _m.Called(ctx, orderID, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.OrderStatus) error); ok {
		r0 = rf(ctx, orderID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrderRepository creates a new instance of OrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderRepository {
	mock := &OrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
