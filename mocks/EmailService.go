// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/G-Villarinho/food-shop-api/models"
)

// EmailService is an autogenerated mock type for the EmailService type
type EmailService struct {
	mock.Mock
}

// SendEmail provides a mock function with given fields: ctx, task
func (_m *EmailService) SendEmail(ctx context.Context, task models.EmailQueueTask) error {
	ret := _m.Called(ctx, task)

	if len(ret) == 0 {
		panic("no return value specified for SendEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EmailQueueTask) error); ok {
		r0 = rf(ctx, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEmailService creates a new instance of EmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmailService(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailService {
	mock := &EmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}