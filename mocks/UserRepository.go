// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/G-Villarinho/food-shop-api/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, ID
func (_m *UserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByID")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.User, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.User); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
