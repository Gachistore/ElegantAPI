package api

import (
	"3legant/types"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func createJWT(account *types.Account) (string, error) {
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

		if claims["userType"] != string(types.UserTypeAdmin) {
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

		if claims["userType"] != string(types.UserTypeRegular) {
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

