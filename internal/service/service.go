package service

import (
	"fmt"
	"log"
	"time"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/model"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(req *model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetSubscription(id string) (*model.Subscription, error)
	UpdateSubscription(id string, req *model.UpdateSubscriptionRequest) error
	DeleteSubscription(id string) error
	ListSubscriptions(limit, offset int) ([]*model.Subscription, error)
	CalculateTotalCost(req *model.CalculateCostRequest) (*model.CalculateCostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	if req.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be positive")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.StartDate == "" {
		return nil, fmt.Errorf("start_date is required")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: must be valid UUID")
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate := startDate.AddDate(0, 1, 0)

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	createdSubscription, err := s.repo.Create(subscription)
	if err != nil {
		return nil, err
	}

	log.Printf("Service: Created subscription %s for user %s", createdSubscription.ID, createdSubscription.UserID)
	return subscription, nil
}

func (s *subscriptionService) GetSubscription(id string) (*model.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetByID(id)
}

func (s *subscriptionService) UpdateSubscription(id string, req *model.UpdateSubscriptionRequest) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	if req.UserID != nil && *req.UserID != "" {
		_, err := uuid.Parse(*req.UserID)
		if err != nil {
			return fmt.Errorf("invalid user_id format: must be valid UUID")
		}
	}

	if req.StartDate != nil && *req.StartDate != "" {
		if _, err := time.Parse("01-2006", *req.StartDate); err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}
	}

	return s.repo.Update(id, req)
}

func (s *subscriptionService) DeleteSubscription(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return s.repo.Delete(id)
}

func (s *subscriptionService) ListSubscriptions(limit, offset int) ([]*model.Subscription, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(limit, offset)
}

func (s *subscriptionService) CalculateTotalCost(req *model.CalculateCostRequest) (*model.CalculateCostResponse, error) {
	if req.StartPeriod == "" || req.EndPeriod == "" {
		return nil, fmt.Errorf("start_period and end_period are required")
	}
	if req.ServiceName == "" {
		return nil, fmt.Errorf("service_name are required")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id are required")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: must be valid UUID")
	}

	total, err := s.repo.CalculateTotalCost(req)
	if err != nil {
		return nil, err
	}

	return &model.CalculateCostResponse{
		TotalCost:   total,
		UserID:      userID,
		ServiceName: req.ServiceName,
		StartPeriod: req.StartPeriod,
		EndPeriod:   req.EndPeriod,
	}, nil
}
