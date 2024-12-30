package services

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestaurantService_CreateRestaurant(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a restaurant successfully", func(t *testing.T) {
		userService := &mocks.UserService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			userService:          userService,
			restaurantRepository: restaurantRepository,
		}

		payload := models.CreateRestaurantPayload{
			RestaurantName: "Test Restaurant",
			Manager: models.CreateUserPayload{
				FullName: "Manager Name",
				Email:    "manager@example.com",
				Phone:    nil,
			},
		}

		managerID := uuid.New()

		userService.On("CreateUser", ctx, payload.Manager, models.Manager).Return(managerID, nil)

		restaurantRepository.On("CreateRestaurant", ctx, mock.MatchedBy(func(r models.Restaurant) bool {
			return r.Name == payload.RestaurantName && r.ManagerID == managerID
		})).Return(nil)

		err := restaurantService.CreateRestaurant(ctx, payload)

		assert.NoError(t, err)
		userService.AssertCalled(t, "CreateUser", ctx, payload.Manager, models.Manager)
		restaurantRepository.AssertCalled(t, "CreateRestaurant", ctx, mock.AnythingOfType("models.Restaurant"))
	})

	t.Run("should return error when user creation fails", func(t *testing.T) {
		userService := &mocks.UserService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			userService:          userService,
			restaurantRepository: restaurantRepository,
		}

		payload := models.CreateRestaurantPayload{
			RestaurantName: "Test Restaurant",
			Manager: models.CreateUserPayload{
				FullName: "Manager Name",
				Email:    "manager@example.com",
				Phone:    nil,
			},
		}

		userService.On("CreateUser", ctx, payload.Manager, models.Manager).Return(uuid.Nil, errors.New("user creation failed"))

		err := restaurantService.CreateRestaurant(ctx, payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user creation failed")
		userService.AssertCalled(t, "CreateUser", ctx, payload.Manager, models.Manager)
		restaurantRepository.AssertNotCalled(t, "CreateRestaurant", ctx, mock.Anything)
	})

	t.Run("should return error when restaurant creation fails", func(t *testing.T) {
		userService := &mocks.UserService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			userService:          userService,
			restaurantRepository: restaurantRepository,
		}

		payload := models.CreateRestaurantPayload{
			RestaurantName: "Test Restaurant",
			Manager: models.CreateUserPayload{
				FullName: "Manager Name",
				Email:    "manager@example.com",
				Phone:    nil,
			},
		}

		managerID := uuid.New()

		userService.On("CreateUser", ctx, payload.Manager, models.Manager).Return(managerID, nil)

		restaurantRepository.On("CreateRestaurant", ctx, mock.MatchedBy(func(r models.Restaurant) bool {
			return r.Name == payload.RestaurantName && r.ManagerID == managerID
		})).Return(errors.New("restaurant creation failed"))

		err := restaurantService.CreateRestaurant(ctx, payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create restaurant")
		userService.AssertCalled(t, "CreateUser", ctx, payload.Manager, models.Manager)
		restaurantRepository.AssertCalled(t, "CreateRestaurant", ctx, mock.AnythingOfType("models.Restaurant"))
	})
}

func TestRestaurantService_CreateOrder(t *testing.T) {
	ctx := context.WithValue(context.Background(), internal.UserIDKey, uuid.New())

	t.Run("should create an order successfully", func(t *testing.T) {
		orderService := &mocks.OrderService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			orderService:         orderService,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateOrderPayload{}
		restaurant := &models.Restaurant{
			BaseModel: models.BaseModel{ID: restaurantID},
		}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(restaurant, nil)

		orderService.On("CreateOrder", ctx, custommerID, restaurantID, payload).Return(nil)

		err := restaurantService.CreateOrder(ctx, restaurantID, payload)

		assert.NoError(t, err)
		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		orderService.AssertCalled(t, "CreateOrder", ctx, custommerID, restaurantID, payload)
	})

	t.Run("should return error when user is not in context", func(t *testing.T) {
		orderService := &mocks.OrderService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			orderService:         orderService,
			restaurantRepository: restaurantRepository,
		}

		invalidCtx := context.Background()
		restaurantID := uuid.New()
		payload := models.CreateOrderPayload{}

		err := restaurantService.CreateOrder(invalidCtx, restaurantID, payload)

		assert.ErrorIs(t, err, models.ErrUserNotFoundInContext)
		restaurantRepository.AssertNotCalled(t, "GetRestaurantByID", invalidCtx, mock.Anything)
		orderService.AssertNotCalled(t, "CreateOrder", invalidCtx, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when restaurant is not found", func(t *testing.T) {
		orderService := &mocks.OrderService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			orderService:         orderService,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateOrderPayload{}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(nil, nil)

		err := restaurantService.CreateOrder(ctx, restaurantID, payload)

		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)
		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		orderService.AssertNotCalled(t, "CreateOrder", ctx, custommerID, restaurantID, payload)
	})

	t.Run("should return error when GetRestaurantByID fails", func(t *testing.T) {
		orderService := &mocks.OrderService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			orderService:         orderService,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateOrderPayload{}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(nil, errors.New("database error"))

		err := restaurantService.CreateOrder(ctx, restaurantID, payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get restaurant by id")
		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		orderService.AssertNotCalled(t, "CreateOrder", ctx, custommerID, restaurantID, payload)
	})

	t.Run("should return error when CreateOrder fails", func(t *testing.T) {
		orderService := &mocks.OrderService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		restaurantService := &restaurantService{
			orderService:         orderService,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateOrderPayload{}
		restaurant := &models.Restaurant{
			BaseModel: models.BaseModel{ID: restaurantID},
		}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(restaurant, nil)

		orderService.On("CreateOrder", ctx, custommerID, restaurantID, payload).Return(errors.New("order creation failed"))

		err := restaurantService.CreateOrder(ctx, restaurantID, payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order creation failed")
		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		orderService.AssertCalled(t, "CreateOrder", ctx, custommerID, restaurantID, payload)
	})
}
