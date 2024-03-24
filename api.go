package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math"
	"net/http"
	"os"
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

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	acc, err := s.store.GetAccountByEmail(req.Email)
	if err != nil {
		return err
	}
	if !acc.ValidPassword(req.Password) {
		fmt.Errorf("not authenticated")
	}
	token, err := createJWT(acc)
	if err != nil {
		return err
	}
	resp := LoginResponse{
		ID:    acc.ID,
		Token: token,
	}

	fmt.Println(acc)
	return WriteJSON(w, http.StatusOK, resp)
}

//ACCOUNT

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}
	//if r.Method == "DELETE" {
	//	return s.handleDeleteAccount(w, r)
	//}

	return nil
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var account Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		return err
	}
	if err := s.store.UpdateAccount(id, &account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	//createCartReq := new(CreateCartRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account, err := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Email,
		createAccountReq.Password, UserTypeRegular)
	//if err := json.NewDecoder(r.Body).Decode(createCartReq); err != nil {
	//	return err
	//}

	if err != nil {
		return err
	}
	id, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}
	account.ID = id
	cart := NewCart(id)
	if err := s.store.CreateCart(cart); err != nil {
		return err
	}
	println(id)
	println(cart.UserID, cart.CartID)
	fmt.Println(account)

	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// PRODUCT

func (s *APIServer) handleProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetProduct(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateProduct(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateProduct(w, r)
	}
	//if r.Method == "DELETE" {
	//	return s.handleDeleteProduct(w, r)
	//}

	return nil
}

func (s *APIServer) handleGetProductByID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	if idStr == "categories" {
		return s.handleGetCategory(w, r)
	}
	if idStr == "reviews" {
		return s.handleGetReview(w, r)
	}
	//if strings.HasPrefix(idStr, "search") {
	//	return s.handleSearchProduct(w, r)
	//}
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		product, err := s.store.GetProductByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, product)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteProduct(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateProduct(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return err
	}
	if err := s.store.UpdateProduct(id, &product); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *APIServer) handleGetProduct(w http.ResponseWriter, r *http.Request) error {
	//product, err := s.store.GetProducts()
	vars := r.URL.Query()
	name := vars.Get("name")
	priceFrom := vars.Get("priceFrom")
	priceTo := vars.Get("priceTo")
	skip := vars.Get("skip")
	limit := vars.Get("limit")
	if priceFrom == "" {
		priceFrom = "0"
	}
	if priceTo == "" {
		priceTo = fmt.Sprintf("%f", math.MaxFloat32)
	}
	params := map[string]any{"name": name, "priceFrom": priceFrom, "priceTo": priceTo, "skip": skip, "limit": limit}
	products, err := s.store.SearchProducts(params)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	createProductReq := new(CreateProductRequest)
	if err := json.NewDecoder(r.Body).Decode(createProductReq); err != nil {
		return err
	}
	product := NewProduct(createProductReq.Name,
		createProductReq.Price,
		createProductReq.Measurements,
		createProductReq.Description,
		createProductReq.Packaging)
	if err := s.store.CreateProduct(product); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	err1, err2, err3 := s.store.DeleteProduct(id)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleGetNewProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetNewProducts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) handleSearchProduct(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()
	name := vars.Get("name")
	priceFrom := vars.Get("priceFrom")
	priceTo := vars.Get("priceTo")
	if priceFrom == "" {
		priceFrom = "0"
	}
	if priceTo == "" {
		priceTo = fmt.Sprintf("%f", math.MaxFloat32)
	}
	params := map[string]any{"name": name, "priceFrom": priceFrom, "priceTo": priceTo}
	products, err := s.store.SearchProducts(params)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

// REVIEW

func (s *APIServer) handleReview(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetReview(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateReview(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateReview(w, r)
	}
	//if r.Method == "DELETE" {
	//	return s.handleDeleteReview(w, r)
	//}

	return nil
}
func (s *APIServer) handleGetReviewByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetReviewByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteReview(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetReview(w http.ResponseWriter, r *http.Request) error {
	reviews, err := s.store.GetReviews()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, reviews)
}

func (s *APIServer) handleUpdateReview(w http.ResponseWriter, r *http.Request) error {
	id, err1 := getID(r)
	if err1 != nil {
		return err1
	}
	var review Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := s.store.UpdateReview(id, &review); err != nil {
		return err
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": id})
}

func (s *APIServer) handleCreateReview(w http.ResponseWriter, r *http.Request) error {
	createReviewReq := new(CreateReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(createReviewReq); err != nil {
		return err
	}

	review := NewReview(createReviewReq.AccID, createReviewReq.ProdID, createReviewReq.RatingGiven, createReviewReq.Text)
	if err := s.store.CreateReview(review); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, review)
}

func (s *APIServer) handleDeleteReview(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteReview(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// CATEGORY

func (s *APIServer) handleGetCategory(w http.ResponseWriter, r *http.Request) error {
	caregories, err := s.store.GetCategories()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, caregories)
}

// CART

func (s *APIServer) HandleCart(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetCart(w, r)
	}
	if r.Method == "POST" {
		return s.handleAddProductToCart(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateProductQuantityInCart(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)

}

func (s *APIServer) handleGetCart(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	prodQuantities := []*ProductQuantity{} // Создаем слайс для хранения пар продукт-количество

	carts, err := s.store.GetCartProductsByUserID(id)
	if err != nil {
		return err
	}

	for _, cart := range carts {
		prod, err := s.store.GetProductByID(cart.ProdID)
		if err != nil {
			return err
		}
		prodQuantity := &ProductQuantity{
			Product:  prod,
			Quantity: cart.Quantity,
		}
		prodQuantities = append(prodQuantities, prodQuantity)
	}

	return WriteJSON(w, http.StatusOK, prodQuantities) // Возвращаем слайс вместо карты
}

func (s *APIServer) handleAddProductToCart(w http.ResponseWriter, r *http.Request) error {
	cartID, err := getID(r)
	vars := r.URL.Query()
	prodID, err := strconv.Atoi(vars.Get("prodID"))
	if err != nil {
		return err
	}
	quantity, err := strconv.Atoi(vars.Get("quantity"))
	if err != nil {
		return err
	}
	err = s.store.AddProductToCart(cartID, prodID, quantity)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"added: ": prodID})
}

func (s *APIServer) handleUpdateProductQuantityInCart(w http.ResponseWriter, r *http.Request) error {
	cartID, err := getID(r)
	if err != nil {
		return err
	}
	vars := r.URL.Query()
	prodID, err := strconv.Atoi(vars.Get("prodID"))
	if err != nil {
		return err
	}
	quantity, err := strconv.Atoi(vars.Get("quantity"))
	if err != nil {
		return err
	}
	if err := s.store.UpdateProductQuantityInCart(cartID, prodID, quantity); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"updated": prodID})
}

// MISC

func (p Product) String() string {
	// Преобразование структуры в JSON-строку
	jsonData, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("Ошибка при преобразовании в JSON: %v", err)
	}
	return string(jsonData)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"accountID": account.ID,
		"userType":  account.UserType,
	}
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(secret))
	return str, err
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func adminMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if claims["userType"] != string(UserTypeAdmin) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func userMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		id, err := getID(r)
		if err != nil {
			permissionDenied(w)
			fmt.Println(id)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if int(claims["accountID"].(float64)) != id {
			permissionDenied(w)
			return
		}

		if claims["userType"] != string(UserTypeRegular) {
			permissionDenied(w)
			return
		}

		fmt.Println(claims)
		handlerFunc(w, r)
	}
}
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected string method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

type APIServer struct {
	listenAddr string
	store      Storage
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

func newAPIServer(listenAddr string, store Storage) *APIServer {
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

func (acc *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(acc.EncryptedPassword), []byte(pw)) == nil
}
