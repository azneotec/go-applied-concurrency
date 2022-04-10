package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/azad/go-applied-concurrency/models"
	"github.com/azad/go-applied-concurrency/repo"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

type handler struct {
	repo repo.Repo
	once sync.Once
}

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	ProductIndex(w http.ResponseWriter, r *http.Request)
	OrderShow(w http.ResponseWriter, r *http.Request)
	OrderInsert(w http.ResponseWriter, r *http.Request)
	Close(w http.ResponseWriter, r *http.Request)
	Stats(w http.ResponseWriter, r *http.Request)
}

func New() (Handler, error) {
	r, err := repo.New()
	if err != nil {
		return nil, err
	}
	h := handler{repo: r}
	return &h, nil
}

// Index returns a simple hello response for the homepage
func (h handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	writeResponse(w, http.StatusOK, "Welcome to the Orders App!", nil)
}

// ProductIndex displays all products in the system
func (h handler) ProductIndex(w http.ResponseWriter, r *http.Request) {
	p := h.repo.GetAllProducts()
	// Send an HTTP status & send the slice
	writeResponse(w, http.StatusOK, p, nil)
}

// OrderShow fetches and displays one selected product
func (h handler) OrderShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]
	// Call the repository method corresponding to the operation
	o, err := h.repo.GetOrder(orderID)
	// Handler any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusNotFound, nil, err)
		return
	}
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, o, nil)
}

// OrderInsert creates a new order with the given parameters
func (h handler) OrderInsert(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid order body:%v", err))
		return
	}
	order, err := h.repo.CreateOrder(item)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	writeResponse(w, http.StatusCreated, order, nil)
}

func (h *handler) Close(w http.ResponseWriter, r *http.Request) {
	h.invokeClose()
	writeResponse(w, http.StatusOK, "The Orders App is now closed!", nil)
}

func (h *handler) invokeClose() {
	h.once.Do(func() {
		h.repo.Close()
	})
}

func (h *handler) Stats(w http.ResponseWriter, r *http.Request) {
	h.invokeClose()
	writeResponse(w, http.StatusOK, "The Orders App is now closed!", nil)
}
