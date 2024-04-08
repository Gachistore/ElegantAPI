package api

import (
	"3legant/storage"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))

	router.HandleFunc("/accounts", adminMiddleware(makeHTTPHandleFunc(s.handleAccount)))
	router.HandleFunc("/accounts/{id}", adminMiddleware(makeHTTPHandleFunc(s.handleGetAccountByID)))

	router.HandleFunc("/products", makeHTTPHandleFunc(s.handleProduct))
	router.HandleFunc("/products/{id}", makeHTTPHandleFunc(s.handleGetProductByID))
	router.HandleFunc("/products/new", makeHTTPHandleFunc(s.handleGetNewProducts))

	router.HandleFunc("/products/reviews", makeHTTPHandleFunc(s.handleReview))
	router.HandleFunc("/products/reviews/{id}", makeHTTPHandleFunc(s.handleGetReviewByID))

	router.HandleFunc("/products/categories", makeHTTPHandleFunc(s.handleGetCategory))
	//router.HandleFunc("/products/search", makeHTTPHandleFunc(s.handleSearchProduct))

	router.HandleFunc("/carts/{id}", userMiddleware(makeHTTPHandleFunc(s.HandleCart)))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// MISC

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

