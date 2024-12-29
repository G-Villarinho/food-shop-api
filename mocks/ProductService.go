// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/G-Villarinho/food-shop-api/models"
	mock "github.com/stretchr/testify/mock"
)

// ProductService is an autogenerated mock type for the ProductService type
type ProductService struct {
	mock.Mock
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
