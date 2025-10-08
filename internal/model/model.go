package model

import (
	"time"

	"github.com/google/uuid"
)

// Subscription представляет подписку пользователя
// @Description Информация о подписке
type Subscription struct {
	ID          int       `json:"id" example:"1"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"1500"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time `json:"start_date" example:"01-2025"`
	EndDate     time.Time `json:"end_date" example:"02-2025"`
}

// CreateSubscriptionRequest представляет запрос на создание подписки
// @Description Тело запроса для создания новой подписки
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus" binding:"required"`
	Price       int    `json:"price" example:"1500" binding:"required,gt=0"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" binding:"required"`
	StartDate   string `json:"start_date" example:"01-2025" binding:"required"`
}

// UpdateSubscriptionRequest представляет запрос на обновление подписки
// @Description Тело запроса для обновления существующей подписки
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Yandex Plus Premium"`
	Price       *int    `json:"price,omitempty" example:"2000"`
	UserID      *string `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   *string `json:"start_date,omitempty" example:"02-2025"`
}

// CalculateCostRequest представляет запрос на расчет стоимости
// @Description Тело запроса для расчета общей стоимости подписок
type CalculateCostRequest struct {
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" binding:"required"`
	ServiceName string `json:"service_name" example:"Yandex Plus" binding:"required"`
	StartPeriod string `json:"start_period" example:"01-2025" binding:"required"`
	EndPeriod   string `json:"end_period" example:"02-2025" binding:"required"`
}

// CalculateCostResponse представляет ответ с расчетом стоимости
// @Description Ответ с результатом расчета общей стоимости
type CalculateCostResponse struct {
	TotalCost   int       `json:"total_cost" example:"18000"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	StartPeriod string    `json:"start_period" example:"01-2025"`
	EndPeriod   string    `json:"end_period" example:"02-2025"`
}
