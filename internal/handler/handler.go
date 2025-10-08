package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/model"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/service"
	// httpSwagger "github.com/swaggo/http-swagger"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling CreateSubscription request")

	var req model.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.CreateSubscription(&req)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling GetSubscription request for ID: %s", id)

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetSubscription(id)
	if err != nil {
		log.Printf("Error getting subscription: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling UpdateSubscription request for ID: %s", id)

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateSubscription(id, &req); err != nil {
		log.Printf("Error updating subscription: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Handling DeleteSubscription request for ID: %s", id)

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSubscription(id); err != nil {
		log.Printf("Error deleting subscription: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling ListSubscriptions request")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	subscriptions, err := h.service.ListSubscriptions(limit, offset)
	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *SubscriptionHandler) CalculateTotalCost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling CalculateTotalCost request")

	var req model.CalculateCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CalculateTotalCost(&req)
	if err != nil {
		log.Printf("Error calculating total cost: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *SubscriptionHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions", h.CreateSubscription)
	mux.HandleFunc("GET /subscriptions", h.GetSubscription)
	mux.HandleFunc("PUT /subscriptions", h.UpdateSubscription)
	mux.HandleFunc("DELETE /subscriptions", h.DeleteSubscription)
	mux.HandleFunc("GET /subscriptions/list", h.ListSubscriptions)
	mux.HandleFunc("POST /subscriptions/total-cost", h.CalculateTotalCost)
	// mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
