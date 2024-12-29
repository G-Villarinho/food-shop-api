// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/G-Villarinho/food-shop-api/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// RestaurantService is an autogenerated mock type for the RestaurantService type
type RestaurantService struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, restaurantID, payload
func (_m *RestaurantService) CreateOrder(ctx context.Context, restaurantID uuid.UUID, payload models.CreateOrderPayload) error {
	ret := _m.Called(ctx, restaurantID, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.CreateOrderPayload) error); ok {
		r0 = rf(ctx, restaurantID, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateRestaurant provides a mock function with given fields: ctx, payload
func (_m *RestaurantService) CreateRestaurant(ctx context.Context, payload models.CreateRestaurantPayload) error {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateRestaurant")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateRestaurantPayload) error); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRestaurantService creates a new instance of RestaurantService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRestaurantService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RestaurantService {
	mock := &RestaurantService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
