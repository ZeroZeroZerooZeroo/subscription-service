package model

import "time"

type Subscription struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
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
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}

type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type ListSubscriptionResponse struct {
	Subscription []SubscriptionResponse `json:"subscriptions"`
	Total        int                    `json:"total"`
	//Limit        int                    `form:"limit,default=50"`
	//Offset       int                    `form:"offset,default=0"`
}

type CalculateCostRequest struct {
	UserID      string `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	ServiceName string `json:"service_name,omitempty"`
	StartPeriod string `json:"start_period" validate:"required"`
	EndPeriod   string `json:"end_period" validate:"required"`
}

type CalculateCostResponse struct {
	TotalCost   int    `json:"total_cost"`
	UserID      string `json:"user_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	StartPeriod string `json:"start_period"`
	EndPeriod   string `json:"end_period"`
}
