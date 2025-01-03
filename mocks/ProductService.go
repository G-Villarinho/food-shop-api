// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/G-Villarinho/food-shop-api/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ProductService is an autogenerated mock type for the ProductService type
type ProductService struct {
	mock.Mock
}

// CreateProduct provides a mock function with given fields: ctx, payload
func (_m *ProductService) CreateProduct(ctx context.Context, payload *models.CreateOrUpdateProductPayload) (*models.Product, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 *models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateOrUpdateProductPayload) (*models.Product, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateOrUpdateProductPayload) *models.Product); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.CreateOrUpdateProductPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProducts provides a mock function with given fields: ctx, productIDs, restaurantID
func (_m *ProductService) DeleteProducts(ctx context.Context, productIDs []uuid.UUID, restaurantID uuid.UUID) error {
	ret := _m.Called(ctx, productIDs, restaurantID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProducts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, productIDs, restaurantID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPopularProducts provides a mock function with given fields: ctx
func (_m *ProductService) GetPopularProducts(ctx context.Context) ([]models.PopularProductResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetPopularProducts")
	}

	var r0 []models.PopularProductResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.PopularProductResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.PopularProductResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.PopularProductResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProduct provides a mock function with given fields: ctx, payload
func (_m *ProductService) UpdateProduct(ctx context.Context, payload *models.CreateOrUpdateProductPayload) (*models.Product, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProduct")
	}

	var r0 *models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateOrUpdateProductPayload) (*models.Product, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateOrUpdateProductPayload) *models.Product); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.CreateOrUpdateProductPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductService creates a new instance of ProductService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductService {
	mock := &ProductService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
