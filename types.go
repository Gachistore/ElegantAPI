package main

import (
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	UserTypeRegular UserType = "Regular"
	UserTypeAdmin   UserType = "Admin"
)

type Account struct {
	ID                int      `json:"id"`
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Email             string   `json:"email"`
	EncryptedPassword string   `json:"-"`
	UserType          UserType `json:"userType"`
}

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Measurements string  `json:"measurements"`
	Description  string  `json:"description"`
	Packaging    string  `json:"packaging"`
}

type Category struct {
	Name string `json:"name"`
}

type Review struct {
	ID          int     `json:"id"`
	AccID       int     `json:"accID"`
	ProdID      int     `json:"prodID"`
	RatingGiven float64 `json:"ratingGiven"`
	Text        string  `json:"text"`
}

type Cart struct {
	CartID int `json:"cartID"`
	UserID int `json:"userID"`
}

type ProductCart struct {
	CartID   int      `json:"cartID"`
	ProdID int `json:"prodID"`
	Quantity int      `json:"quantity"`
}

type ProductQuantity struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

func NewAccount(firstName, lastName, email, password string, userType UserType) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		EncryptedPassword: string(encpw),
		UserType:          userType,
	}, nil
}

func NewReview(accId, prodID int, ratingGiven float64, text string) *Review {
	return &Review{
		AccID:       accId,
		ProdID:      prodID,
		RatingGiven: ratingGiven,
		Text:        text,
	}
}

func NewProduct(name string, price float64, measurements, description, packaging string) *Product {
	return &Product{
		Name:         name,
		Price:        price,
		Measurements: measurements,
		Description:  description,
		Packaging:    packaging,
	}
}

func NewCart(userID int) *Cart {
	return &Cart{
		UserID: userID,
	}
}

type CreateAccountRequest struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	UserType  UserType `json:"userType"`
}

type CreateProductRequest struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Measurements string  `json:"measurements"`
	Description  string  `json:"description"`
	Packaging    string  `json:"packaging"`
}

type CreateReviewRequest struct {
	AccID       int     `json:"accID"`
	ProdID      int     `json:"prodID"`
	RatingGiven float64 `json:"ratingGiven"`
	Text        string  `json:"text"`
}

type CreateCartRequest struct {
	UserID int `json:"userID"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}
