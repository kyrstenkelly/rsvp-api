package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

/* From tutorial: https://auth0.com/blog/authentication-in-golang/ */

// SigningKey global signing key
var SigningKey = []byte("secret")

// GetTokenHandler handles getting a JWT
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = true
	claims["name"] = "Ado Kukic" // TODO: Update
	claims["exp"] = time.Now().Add(time.Hour * 7 * 24).Unix()

	tokenString, _ := token.SignedString(SigningKey)

	w.Write([]byte(tokenString))
}
