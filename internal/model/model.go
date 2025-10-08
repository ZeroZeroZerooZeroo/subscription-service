package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	UserID      *string `json:"user_id"`
	StartDate   *string `json:"start_date"`
}

type CalculateCostRequest struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	StartPeriod string `json:"start_period"`
	EndPeriod   string `json:"end_period"`
}

type CalculateCostResponse struct {
	TotalCost   int       `json:"total_cost"`
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	StartPeriod string    `json:"start_period"`
	EndPeriod   string    `json:"end_period"`
}
