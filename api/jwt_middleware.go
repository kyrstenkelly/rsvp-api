package api

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/kyrstenkelly/rsvp-api/handlers"
)

var authenticatedEndpoints = []string{"/admin"}

// JWTHeader constant
const JWTHeader = "Authorization"

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return handlers.SigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
