package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
)

type CustomJWTClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

var JwtSecret []byte = []byte("thepolyglotdeveloper")

func ValidateJWT(t string) (interface{}, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return nil, errors.New(`{ "message": "` + err.Error() + `"}`)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var tokenData CustomJWTClaims
		mapstructure.Decode(claims, &tokenData)
		return tokenData, nil
	} else {
		return nil, errors.New(`{ "message": "invalid token" }`)
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authorizeHeader := request.Header.Get("authorization")
		if authorizeHeader != "" {
			bearerToken := strings.Split(authorizeHeader, " ")
			if len(bearerToken) == 2 {
				decoded, err := ValidateJWT(bearerToken[1])
				if err != nil {
					response.Header().Add("content-type", "application/json")
					response.WriteHeader(500)
					response.Write([]byte(err.Error()))
					return
				}
				context.Set(request, "decoded", decoded)
				next(response, request)
			}
		} else {
			response.Header().Add("content-type", "application/json")
			response.WriteHeader(500)
			response.Write([]byte(`{ "message": "header required" }`))
		}
	})
}
