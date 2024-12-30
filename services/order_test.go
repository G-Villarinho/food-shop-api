package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("should create order successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		productRepository := &mocks.ProductRepository{}
		orderItemService := &mocks.OrderItemService{}

		orderService := &orderService{
			orderRepository:   orderRepository,
			productRepository: productRepository,
			orderItemService:  orderItemService,
		}

		custommerID := uuid.New()
		restaurantID := uuid.New()

		product1ID := uuid.New()
		product2ID := uuid.New()
		products := []models.Product{
			{BaseModel: models.BaseModel{ID: product1ID}, PriceInCents: 1000},
			{BaseModel: models.BaseModel{ID: product2ID}, PriceInCents: 2000},
		}

		items := []models.CreateOrderItemPayload{
			{ProductID: product1ID, Quantity: 2},
			{ProductID: product2ID, Quantity: 1},
		}

		payload := models.CreateOrderPayload{
			Items: items,
		}

		productRepository.On("GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID).Return(products, nil)

		orderItems := []models.OrderItem{
			{ProductID: product1ID, Quantity: 2, PriceInCents: 1000},
			{ProductID: product2ID, Quantity: 1, PriceInCents: 2000},
		}

		orderItemService.On("ValidateAndCalculateOrderItems", mock.Anything, products, items).Return(&models.OrderItemSummary{
			OrderItems:   orderItems,
			TotalInCents: 4000,
		}, nil)

		orderRepository.On("CreateOrderWithItems", mock.Anything, mock.Anything, orderItems).Return(nil)

		err := orderService.CreateOrder(context.Background(), custommerID, restaurantID, payload)

		assert.NoError(t, err)

		productRepository.AssertCalled(t, "GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID)
		orderItemService.AssertCalled(t, "ValidateAndCalculateOrderItems", mock.Anything, products, items)
		orderRepository.AssertCalled(t, "CreateOrderWithItems", mock.Anything, mock.Anything, orderItems)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		productRepository := &mocks.ProductRepository{}
		orderItemService := &mocks.OrderItemService{}

		orderService := &orderService{
			orderRepository:   orderRepository,
			productRepository: productRepository,
			orderItemService:  orderItemService,
		}

		custommerID := uuid.New()
		restaurantID := uuid.New()

		productRepository.On("GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID).Return(nil, fmt.Errorf("products not found"))

		items := []models.CreateOrderItemPayload{
			{ProductID: uuid.New(), Quantity: 1},
		}

		payload := models.CreateOrderPayload{
			Items: items,
		}

		err := orderService.CreateOrder(context.Background(), custommerID, restaurantID, payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get products by ids and restaurant id")

		productRepository.AssertCalled(t, "GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID)
		orderItemService.AssertNotCalled(t, "ValidateAndCalculateOrderItems", mock.Anything, mock.Anything, mock.Anything)
		orderRepository.AssertNotCalled(t, "CreateOrderWithItems", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when validation of items fails", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		productRepository := &mocks.ProductRepository{}
		orderItemService := &mocks.OrderItemService{}

		orderService := &orderService{
			orderRepository:   orderRepository,
			productRepository: productRepository,
			orderItemService:  orderItemService,
		}

		custommerID := uuid.New()
		restaurantID := uuid.New()

		product1ID := uuid.New()
		products := []models.Product{
			{BaseModel: models.BaseModel{ID: product1ID}, PriceInCents: 1000},
		}

		items := []models.CreateOrderItemPayload{
			{ProductID: product1ID, Quantity: 1},
		}

		payload := models.CreateOrderPayload{
			Items: items,
		}

		productRepository.On("GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID).Return(products, nil)

		orderItemService.On("ValidateAndCalculateOrderItems", mock.Anything, products, items).Return(nil, fmt.Errorf("validation failed"))

		err := orderService.CreateOrder(context.Background(), custommerID, restaurantID, payload)

		assert.Error(t, err)
		assert.Equal(t, models.ErrSomeProductsNotFound, err)

		productRepository.AssertCalled(t, "GetProductsByIDsAndRestaurantID", mock.Anything, mock.Anything, restaurantID)
		orderItemService.AssertCalled(t, "ValidateAndCalculateOrderItems", mock.Anything, products, items)
		orderRepository.AssertNotCalled(t, "CreateOrderWithItems", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_GetPaginatedOrdersByRestaurantID(t *testing.T) {
	t.Run("should return paginated orders successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}

		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		pagination := &models.OrderPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		mockOrders := []models.Order{
			{
				BaseModel:    models.BaseModel{ID: uuid.New()},
				Custommer:    models.User{FullName: "Customer 1"},
				Status:       models.Pending,
				TotalInCents: 1000,
			},
			{
				BaseModel:    models.BaseModel{ID: uuid.New()},
				Custommer:    models.User{FullName: "Customer 2"},
				Status:       models.Delivered,
				TotalInCents: 2000,
			},
		}

		mockPaginatedOrders := &models.PaginatedResponse[models.Order]{
			Data:       mockOrders,
			Total:      2,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
		}

		orderRepository.On("GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination).Return(mockPaginatedOrders, nil)

		response, err := orderService.GetPaginatedOrdersByRestaurantID(ctx, pagination)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.TotalPages)
		assert.Equal(t, "Customer 1", response.Data[0].CustommerName)
		assert.Equal(t, "Customer 2", response.Data[1].CustommerName)

		orderRepository.AssertCalled(t, "GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}

		orderService := &orderService{
			orderRepository: orderRepository,
		}

		invalidCtx := context.Background()

		pagination := &models.OrderPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		response, err := orderService.GetPaginatedOrdersByRestaurantID(invalidCtx, pagination)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		orderRepository.AssertNotCalled(t, "GetPaginatedOrdersByRestaurantID", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}

		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		pagination := &models.OrderPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		orderRepository.On("GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination).Return(nil, fmt.Errorf("database error"))

		response, err := orderService.GetPaginatedOrdersByRestaurantID(ctx, pagination)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "get paginated orders by restaurant ID")

		orderRepository.AssertCalled(t, "GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination)
	})

	t.Run("should return nil when no orders are found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}

		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		pagination := &models.OrderPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		orderRepository.On("GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination).Return(nil, nil)

		response, err := orderService.GetPaginatedOrdersByRestaurantID(ctx, pagination)

		assert.NoError(t, err)
		assert.Nil(t, response)

		orderRepository.AssertCalled(t, "GetPaginatedOrdersByRestaurantID", ctx, restaurantID, pagination)
	})
}

func TestOrderService_CancelOrder(t *testing.T) {
	t.Run("should cancel order successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Canceled).Return(nil)

		err := orderService.CancelOrder(ctx, orderID)

		assert.NoError(t, err)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Canceled)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		invalidCtx := context.Background()
		orderID := uuid.New()

		err := orderService.CancelOrder(invalidCtx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		orderRepository.AssertNotCalled(t, "GetOrderByID", invalidCtx, mock.Anything, mock.Anything)
		orderRepository.AssertNotCalled(t, "UpdateStatus", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order is not found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(nil, nil)

		err := orderService.CancelOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderNotFound)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order does not belong to restaurant", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: uuid.New(), // Different restaurant ID
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.CancelOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderDoesNotBelongToRestaurant)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order cannot be canceled", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Delivered, // Status not eligible for cancellation
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.CancelOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderCannotBeCancelled)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to update status", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Canceled).Return(fmt.Errorf("database error"))

		err := orderService.CancelOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update status")

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Canceled)
	})
}

func TestOrderService_ApproveOrder(t *testing.T) {
	t.Run("should approve order successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Processing).Return(nil)

		err := orderService.ApproveOrder(ctx, orderID)

		assert.NoError(t, err)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Processing)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		invalidCtx := context.Background()
		orderID := uuid.New()

		err := orderService.ApproveOrder(invalidCtx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		orderRepository.AssertNotCalled(t, "GetOrderByID", invalidCtx, mock.Anything, mock.Anything)
		orderRepository.AssertNotCalled(t, "UpdateStatus", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order is not found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(nil, nil)

		err := orderService.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderNotFound)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order does not belong to restaurant", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: uuid.New(), // Different restaurant ID
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderDoesNotBelongToRestaurant)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order cannot be approved", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Delivered, // Status not eligible for approval
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrOrderCannotBeApproved)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to update status", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Processing).Return(fmt.Errorf("database error"))

		err := orderService.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update status")

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Processing)
	})
}

func TestOrderService_DispatchOrder(t *testing.T) {
	t.Run("should dispatch order successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Processing,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Delivering).Return(nil)

		err := orderService.DispatchOrder(ctx, orderID)

		assert.NoError(t, err)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Delivering)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		invalidCtx := context.Background()
		orderID := uuid.New()

		err := orderService.DispatchOrder(invalidCtx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		orderRepository.AssertNotCalled(t, "GetOrderByID", invalidCtx, mock.Anything, mock.Anything)
		orderRepository.AssertNotCalled(t, "UpdateStatus", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order is not found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(nil, nil)

		err := orderService.DispatchOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderNotFound)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order does not belong to restaurant", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: uuid.New(),
			Status:       models.Processing,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.DispatchOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderDoesNotBelongToRestaurant)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order cannot be dispatched", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Pending,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.DispatchOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrOrderCannotBeDispatched)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to update status", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Processing,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Delivering).Return(fmt.Errorf("database error"))

		err := orderService.DispatchOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update status")

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Delivering)
	})
}

func TestOrderService_DeliverOrder(t *testing.T) {
	t.Run("should deliver order successfully", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Delivering,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Delivered).Return(nil)

		err := orderService.DeliverOrder(ctx, orderID)

		assert.NoError(t, err)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Delivered)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		invalidCtx := context.Background()
		orderID := uuid.New()

		err := orderService.DeliverOrder(invalidCtx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		orderRepository.AssertNotCalled(t, "GetOrderByID", invalidCtx, mock.Anything, mock.Anything)
		orderRepository.AssertNotCalled(t, "UpdateStatus", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order is not found", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(nil, nil)

		err := orderService.DeliverOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderNotFound)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order does not belong to restaurant", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: uuid.New(),
			Status:       models.Delivering,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.DeliverOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorOrderDoesNotBelongToRestaurant)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when order cannot be delivered", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Processing,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)

		err := orderService.DeliverOrder(ctx, orderID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrOrderCannotBeDelivered)

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertNotCalled(t, "UpdateStatus", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to update status", func(t *testing.T) {
		orderRepository := &mocks.OrderRepository{}
		orderService := &orderService{
			orderRepository: orderRepository,
		}

		restaurantID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

		orderID := uuid.New()
		mockOrder := &models.Order{
			BaseModel:    models.BaseModel{ID: orderID},
			RestaurantID: restaurantID,
			Status:       models.Delivering,
		}

		orderRepository.On("GetOrderByID", ctx, orderID, false).Return(mockOrder, nil)
		orderRepository.On("UpdateStatus", ctx, orderID, models.Delivered).Return(fmt.Errorf("database error"))

		err := orderService.DeliverOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update status")

		orderRepository.AssertCalled(t, "GetOrderByID", ctx, orderID, false)
		orderRepository.AssertCalled(t, "UpdateStatus", ctx, orderID, models.Delivered)
	})
}
