package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) (*model.Subscription, error)
	GetByID(id string) (*model.Subscription, error)
	Update(id string, req *model.UpdateSubscriptionRequest) error
	Delete(id string) error
	List(limit, offset int) ([]*model.Subscription, error)
	CalculateTotalCost(req *model.CalculateCostRequest) (int, error)
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *model.Subscription) (*model.Subscription, error) {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var createdID int
	err := r.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&createdID)

	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	sub.ID = createdID
	log.Printf("Subscription created successfully: %s", sub.ID)
	return sub, nil
}

func (r *subscriptionRepo) GetByID(id string) (*model.Subscription, error) {

	query := `SELECT id, service_name, price, user_id, start_date, end_date
    FROM subscriptions WHERE id=$1`

	var sub model.Subscription

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: must be integer")
	}

	err = r.db.QueryRow(query, idInt).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}

	if err != nil {
		log.Printf("Error getting subscription by ID: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	log.Printf("Subscription retrieved: %s", id)
	return &sub, nil
}

func (r *subscriptionRepo) Update(id string, req *model.UpdateSubscriptionRequest) error {
	currentSub, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get current subscription: %w", err)
	}

	log.Printf("Current subscription: ID=%d, EndDate=%v", currentSub.ID, currentSub.EndDate)

	query := `UPDATE subscriptions 
	SET service_name = COALESCE($1, service_name),price = COALESCE($2, price),
    user_id = COALESCE($3, user_id),start_date = COALESCE($4, start_date),
    end_date = COALESCE($5, end_date) WHERE id = $6`

	var startDate, endDate interface{}

	if req.StartDate != nil && *req.StartDate != "" {
		parsedDate, err := time.Parse("01-2006", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 1, 0)
		log.Printf("New dates - start: %v, end: %v", startDate, endDate)
	} else {
		startDate = nil
		endDate = nil
	}

	var serviceName, price, userID interface{}

	if req.ServiceName != nil {
		serviceName = *req.ServiceName
	} else {
		serviceName = nil
	}

	if req.Price != nil {
		price = *req.Price
	} else {
		price = nil
	}

	if req.UserID != nil && *req.UserID != "" {
		if _, err := uuid.Parse(*req.UserID); err != nil {
			return fmt.Errorf("invalid user_id format")
		}
		userID = *req.UserID
	} else {
		userID = nil
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id format: must be integer")
	}

	log.Printf("Executing update: service=%v, price=%v, user=%v, start=%v, end=%v",
		serviceName, price, userID, startDate, endDate)

	result, err := r.db.Exec(query, serviceName, price, userID, startDate, endDate, idInt)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Subscription updated successfully: %s", id)
	return nil
}

func (r *subscriptionRepo) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id=$1`

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id format: must be integer")
	}

	result, err := r.db.Exec(query, idInt)
	if err != nil {
		log.Printf("Error deleting subscription: %v", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rewsAffected, _ := result.RowsAffected()

	if rewsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Subscription deleted: %s", id)
	return nil
}

func (r *subscriptionRepo) List(limit, offset int) ([]*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)

	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	defer rows.Close()

	var subscriptions []*model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, &sub)
	}

	log.Printf("Listed %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (r *subscriptionRepo) CalculateTotalCost(req *model.CalculateCostRequest) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
    WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)
    AND user_id = $3 AND service_name = $4`

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format: %w", err)
	}

	startPeriod, err := time.Parse("01-2006", req.StartPeriod)
	if err != nil {
		return 0, fmt.Errorf("invalid start_period format: %w", err)
	}

	endPeriod, err := time.Parse("01-2006", req.EndPeriod)

	if err != nil {
		return 0, fmt.Errorf("invalid end_period format: %w", err)
	}

	var totalCost int
	err = r.db.QueryRow(query, endPeriod, startPeriod, userID, req.ServiceName).Scan(&totalCost)

	if err != nil {
		log.Printf("Error calculating total cost: %v", err)
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	log.Printf("Total cost calculated: %d for period %s to %s", totalCost, req.StartPeriod, req.EndPeriod)
	return totalCost, nil

}
