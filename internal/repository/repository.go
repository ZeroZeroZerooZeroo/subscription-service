package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
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

func (r *subscriptionRepo) Create(sub *model.Subscription) error {
	query := `INSERT INTO subscriptions (id,service_name,price,user_id,start_date,end_date)
	VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := r.db.Exec(query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)

	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	log.Printf("Subscription created successfully: %s", sub.ID)
	return nil
}

func (r *subscriptionRepo) GetByID(id string) (*model.Subscription, error) {
	query := `SELECT id,service_name,price,user_id,start_date,end_date
	FROM subscriptions WHERE id=$1`

	var sub model.Subscription
	var endDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}

	if err != nil {
		log.Printf("Error getting subscription by ID: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if endDate.Valid {
		sub.EndDate = endDate.Time
	}

	log.Printf("Subscription retrieved: %s", id)
	return &sub, nil
}

func (r *subscriptionRepo) Update(id string, req *model.UpdateSubscriptionRequest) error {
	query := `UPDATE subscriptions SET 
	service_name = COALESCE($1, service_name),price = COALESCE($2, price),
	user_id = COALESCE($3, user_id),start_date = COALESCE($4, start_date),
	end_date = $5 WHERE id = $6`

	var startDate interface{}

	if req.StartDate != "" {
		parsedDate, err := time.Parse("01-2006", req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}
		startDate = parsedDate
	} else {
		startDate = nil
	}

	var endDate interface{}

	result, err := r.db.Exec(query, req.ServiceName, req.Price, req.UserID, startDate, endDate, id)

	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)

	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Subscription updated: %s", id)
	return nil
}
func (r *subscriptionRepo) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id=$1`

	result, err := r.db.Exec(query, id)
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
	query := `SELECT id,service_name,price,user_id,start_date,end_date
	FROM subscriptions ORDER BY start_date DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)

	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	defer rows.Close()

	var subscriptions []*model.Subscription

	for rows.Next() {
		var sub model.Subscription
		var endDate sql.NullTime

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		if endDate.Valid {
			sub.EndDate = endDate.Time
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

	startPeriod, err := time.Parse("01-2006", req.StartPeriod)
	if err != nil {
		return 0, fmt.Errorf("invalid start_period format: %w", err)
	}

	endPeriod, err := time.Parse("01-2006", req.EndPeriod)

	if err != nil {
		return 0, fmt.Errorf("invalid end_period format: %w", err)
	}

	var totalCost int
	err = r.db.QueryRow(query, endPeriod, startPeriod, req.UserID, req.ServiceName).Scan(&totalCost)

	if err != nil {
		log.Printf("Error calculating total cost: %v", err)
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	log.Printf("Total cost calculated: %d for period %s to %s", totalCost, req.StartPeriod, req.EndPeriod)
	return totalCost, nil

}
