package main

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Measurements string  `json:"measurements"`
	Description  string  `json:"description"`
	Packaging    string  `json:"packaging"`
}

type Review struct {
	ID          int     `json:"id"`
	AccID       int     `json:"accID"`
	ProdID      int     `json:"prodID"`
	RatingGiven float64 `json:"ratingGiven"`
	Text        string  `json:"text"`
}

func NewAccount(firstName, lastName, email string) *Account {
	return &Account{FirstName: firstName,
		LastName: lastName,
		Email:    email}
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

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
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
