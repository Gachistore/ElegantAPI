package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(int, *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByEmail(string) (*Account, error)
	GetAccountByID(int) (*Account, error)

	CreateProduct(*Product) error
	DeleteProduct(int) (error, error, error)
	UpdateProduct(int, *Product) error
	GetProducts() ([]*Product, error)
	GetProductByID(int) (*Product, error)
	GetNewProducts() ([]*Product, error)
	SearchProducts(map[string]any) ([]*Product, error)

	CreateReview(*Review) error
	DeleteReview(int) error
	UpdateReview(int, *Review) error
	GetReviews() ([]*Review, error)
	GetReviewByID(int) (*Review, error)

	CreateCart(*Cart) error
	UpdateProductQuantityInCart(int, int, int) error
	DeleteProductFromCart(int, int) error
	GetCartProductsByUserID(int) ([]*ProductCart, error)
	AddProductToCart(int, int, int) error

	GetCategories() ([]*Category, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=balls sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil

}

func (s *PostgresStore) Init() []error {
	var errors []error
	errors = append(errors, s.CreateAccountTable())
	errors = append(errors, s.CreateProductTable())
	errors = append(errors, s.CreateReviewTable())
	errors = append(errors, s.CreateCategoryTable())
	errors = append(errors, s.CreateProductCategoryTable())
	errors = append(errors, s.CreateProductReviewTable())
	errors = append(errors, s.CreateCartTable())
	errors = append(errors, s.CreateCartProductTable())
	return errors
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account( 
			id serial primary key,			
			first_name varchar(50),
            last_name varchar(50),
            e_mail varchar(50),
    		encrypted_password varchar(100),
    		user_type varchar(50)
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) (int, error) {
	query := `insert into account (first_name, last_name, e_mail, encrypted_password, user_type)
								   values ($1, $2, $3, $4, $5) returning id`
	var id int
	err := s.db.QueryRow(query,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.EncryptedPassword,
		acc.UserType,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	rows, err := s.db.Query(`select * from account where e_mail = $1`, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with email %s not found", email)
}

func (s *PostgresStore) UpdateAccount(id int, account *Account) error {
	_, err := s.db.Query(`UPDATE account SET first_name=$2, last_name=$3, e_mail=$4 WHERE id=$1`,
		id, account.FirstName, account.LastName, account.Email)
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	cart, err := s.getCartByUserID(id)
	cartID := cart.CartID
	if err != nil {
		return err
	}
	_, err = s.db.Query(`delete from cart_product where cart_id = $1`, cartID)
	if err != nil {
		return err
	}
	_, err = s.db.Query(`delete from cart where id = $1`, cartID)
	if err != nil {
		return err
	}
	_, err = s.db.Query(`delete from account where id = $1`, id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`select * from account where id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.EncryptedPassword,
		&account.UserType,
	)
	return account, err
}

// PRODUCT

func (s *PostgresStore) CreateProductTable() error {
	query := `create table if not exists product( 
    		id serial primary key,
			name         varchar(50),   
			price        real, 
			measurements varchar(50),
			description  varchar(500),
			packaging    varchar(50)			
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateProduct(product *Product) error {

	query := `insert into product
    		(name, price, measurements, description, packaging)
								   values ($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(query,
		product.Name,
		product.Price,
		product.Measurements,
		product.Description,
		product.Packaging)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func scanIntoProduct(rows *sql.Rows) (*Product, error) {
	product := new(Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Measurements,
		&product.Description,
		&product.Packaging)
	return product, err
}

func (s *PostgresStore) UpdateProduct(id int, product *Product) error {
	_, err := s.db.Query(`update product set name=$2, price=$3, measurements=$4, description=$5, packaging=$6 WHERE id=$1`,
		id, product.Name, product.Price, product.Measurements, product.Description, product.Packaging)
	return err
}

func (s *PostgresStore) DeleteProduct(id int) (error, error, error) {
	_, err1 := s.db.Query(`delete from product_review where prodid = $1`, id)
	_, err2 := s.db.Query(`delete from product_category where prodid = $1`, id)
	_, err3 := s.db.Query(`delete from product where id = $1`, id)
	return err1, err2, err3
}

func (s *PostgresStore) GetProductByID(id int) (*Product, error) {
	rows, err := s.db.Query(`select * from product where id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoProduct(rows)
	}
	return nil, fmt.Errorf("product %d not found", id)
}

func (s *PostgresStore) GetProducts() ([]*Product, error) {
	rows, err := s.db.Query(`select * from product`)
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *PostgresStore) GetNewProducts() ([]*Product, error) {
	rows, err := s.db.Query(`select * from product order by id desc limit 5`)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *PostgresStore) SearchProducts(params map[string]any) ([]*Product, error) {
	name := AnyToStr(params["name"])
	name = "%" + name + "%"
	priceFrom := params["priceFrom"]
	priceTo := params["priceTo"]
	skip := params["skip"]
	limit := params["limit"]
	if skip == "" {
		skip = "1"
	}
	if limit == "" {
		limit = "2147483647"
	}
	rows, err := s.db.Query(`select * from product where lower(name) like lower($1) and  price >= $2 and price <=$3 order by id offset $4 - 1 limit $5 - $4 + 1 `,
		name, priceFrom, priceTo, skip, limit)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// REVIEW

func (s *PostgresStore) CreateReviewTable() error {
	query := `create table if not exists review( 
    			id serial primary key,
				accID serial references account(id),
    			prodID serial references product(id),
				rating_given real,
				text varchar(200)
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateReview(rev *Review) error {
	query := `insert into review (accID, prodID, rating_given, text)
								   values ($1, $2, $3, $4)`
	resp, err := s.db.Query(query,
		rev.AccID,
		rev.ProdID,
		rev.RatingGiven,
		rev.Text)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) UpdateReview(id int, review *Review) error {
	_, err := s.db.Query(`UPDATE review SET rating_given=$2, text=$3 WHERE id=$1`,
		id, review.RatingGiven, review.Text)
	return err
}

func (s *PostgresStore) DeleteReview(id int) error {
	//_, err1 := s.db.Query(`delete from product_review where revewid = $1`, id)
	_, err := s.db.Query(`delete from review where id = $1`, id)
	return err
}

func (s *PostgresStore) GetReviewByID(id int) (*Review, error) {
	rows, err := s.db.Query(`select * from review where id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoReview(rows)
	}
	return nil, fmt.Errorf("review %d not found", id)
}

func (s *PostgresStore) GetReviews() ([]*Review, error) {
	rows, err := s.db.Query(`select * from review`)
	if err != nil {
		return nil, err
	}
	reviews := []*Review{}
	for rows.Next() {
		review, err := scanIntoReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func scanIntoReview(rows *sql.Rows) (*Review, error) {
	review := new(Review)
	err := rows.Scan(
		&review.ID,
		&review.AccID,
		&review.ProdID,
		&review.RatingGiven,
		&review.Text)
	return review, err
}

// CATEGORY

func (s *PostgresStore) GetCategories() ([]*Category, error) {
	rows, err := s.db.Query(`select * from category`)
	if err != nil {
		return nil, err
	}
	categories := []*Category{}
	for rows.Next() {
		category, err := scanIntoCategory(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func scanIntoCategory(rows *sql.Rows) (*Category, error) {
	category := new(Category)
	err := rows.Scan(
		&category.Name)
	return category, err
}

//CART

func (s *PostgresStore) CreateCart(cart *Cart) error {
	query := `insert into cart (user_id) values ($1)`
	resp, err := s.db.Query(query,
		cart.UserID,
	)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) AddProductToCart(userID, prodID, quantity int) error {
	cart, err := s.getCartByUserID(userID)
	cartID := cart.CartID
	if err != nil {
		return err
	}
	query := `insert into cart_product (cart_id, product_id, quantity)
								   values ($1, $2, $3)`
	resp, err := s.db.Query(query,
		cartID,
		prodID,
		quantity,
	)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil

}

func (s *PostgresStore) UpdateProductQuantityInCart(cartID, prodID, quantity int) error {
	_, err := s.db.Query(`UPDATE cart_product SET quantity=$3 WHERE cart_id=$1 and product_id=$2 `,
		cartID, prodID, quantity)
	return err
}

func (s *PostgresStore) DeleteProductFromCart(cartID, productID int) error {
	_, err := s.db.Query(`delete from cart_product where cart_id = $1 and product_id = $2`, cartID, productID)
	return err
}

func (s *PostgresStore) getCartByUserID(userID int) (*Cart, error) {
	rows, err := s.db.Query(`select * from cart where user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	cart := new(Cart)
	for rows.Next() {
		cart, err = scanIntoCart(rows)
		if err != nil {
			return nil, err
		}
	}
	return cart, nil
}

func (s *PostgresStore) GetCartProductsByUserID(userID int) ([]*ProductCart, error) {
	cart, err := s.getCartByUserID(userID)
	var products []*ProductCart
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`select * from cart_product where cart_id = $1`, cart.CartID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product, err := s.scanIntoCartProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if products != nil {
		return products, nil
	}
	return nil,  fmt.Errorf("%v not found", userID)
}


func scanIntoCart(rows *sql.Rows) (*Cart, error) {
	cart := new(Cart)
	err := rows.Scan(
		&cart.CartID,
		&cart.UserID,
	)
	return cart, err
}

func (s *PostgresStore) scanIntoCartProduct(rows *sql.Rows) (*ProductCart, error) {
	prodCart := new(ProductCart)
	//prod := new(*Product)
	//prod, err := s.GetProductByID(prodID)
	err := rows.Scan(
		&prodCart.CartID,
		&prodCart.ProdID,
		&prodCart.Quantity,
	)
	if err != nil {
		return nil, err
	}
	return prodCart, nil
}

// MISC

func (s *PostgresStore) CreateProductCategoryTable() error {
	query := `create table if not exists product_category( 
			prodID serial references product(id),
			category_name varchar(50) references category(name),
			constraint product_category_pk primary key (prodID, category_name)
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateCategoryTable() error {
	query := `create table if not exists category( 
			name varchar(50) primary key
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateProductReviewTable() error {
	query := `create table if not exists product_review( 
			prodID serial references product(id),
			reviewID serial references review(id),
			constraint product_review_pk primary key (prodID, reviewID)
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateCartTable() error {
	query := `create table if not exists cart( 
			id serial primary key,
			user_id serial references account(id)			
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateCartProductTable() error {
	query := `create table if not exists cart_product( 
			cart_id serial references cart(id),
			product_id serial references product(id),
    		quantity serial, 
			constraint cart_product_pk primary key (cart_id, product_id)
		)`

	_, err := s.db.Exec(query)
	return err
}

func AnyToStr(param any) string {
	var str string
	switch v := param.(type) {
	case string:
		str = v
	default:
		str = fmt.Sprintf("%v", param)
	}
	return str
}
