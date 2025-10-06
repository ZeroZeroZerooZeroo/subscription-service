package repository

import (
	"database/sql"
	"fmt"
	"log"

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
	query := `INSERT INTO subscritions (id,service_name,price,user_id,start_date,end_date)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := r.db.Exec(query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)

	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	log.Printf("Subscription created successfully: %s", sub.ID)
	return nil
}

func (r *subscriptionRepo) GetByID(id string)(*model.Subscription,error){

}

func (r *subscriptionRepo) Update(id string, req *model.UpdateSubscriptionRequest) error{

}
func (r *subscriptionRepo)Delete(id string) error{

}
func (r *subscriptionRepo)List(limit, offset int) ([]*model.Subscription, error){

}
func (r *subscriptionRepo)CalculateTotalCost(req *model.CalculateCostRequest) (int, error){
	
}
