package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type Author struct {
	Id        string `json:"id,omitempty" validate:"omitempty,uuid"`
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	Username  string `json:"username,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,gte=4"`
}

func RegisterEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var author Author
	json.NewDecoder(request.Body).Decode(&author)
	validate := validator.New()
	err := validate.Struct(author)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(author.Password), 10)
	author.Id = uuid.Must(uuid.NewV4()).String()
	author.Password = string(hash)
	authors = append(authors, author)
	json.NewEncoder(response).Encode(authors)
}

func LoginEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var data Author
	json.NewDecoder(request.Body).Decode(&data)
	validate := validator.New()
	err := validate.StructExcept(data, "Firstname", "Lastname")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for _, author := range authors {
		if author.Username == data.Username {
			err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(data.Password))
			if err != nil {
				response.WriteHeader(500)
				response.Write([]byte(`{ "message": "invalid password" }`))
				return
			}

			claims := CustomJWTClaims{
				Id: author.Id,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
					Issuer:    "The Polyglot Developer",
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString(JwtSecret)
			response.Write([]byte(`{ "token": "` + tokenString + `" }`))
			return
		}
	}
	response.Write([]byte(`{ "message": "invalid username" }`))

}

func AuthorRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(authors)
}

func AuthorRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for _, author := range authors {
		if author.Id == params["id"] {
			json.NewEncoder(response).Encode(author)
			return
		}
	}
	json.NewEncoder(response).Encode(Author{})
}

func AuthorUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var changes Author
	params := mux.Vars(request)
	json.NewDecoder(request.Body).Decode(&changes)
	validate := validator.New()
	err := validate.StructExcept(changes, "Firstname", "Lastname", "Username", "Password")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for index, author := range authors {
		if author.Id == params["id"] {
			if changes.Firstname != "" {
				author.Firstname = changes.Firstname
			}
			if changes.Lastname != "" {
				author.Lastname = changes.Lastname
			}
			if changes.Username != "" {
				author.Username = changes.Username
			}
			if changes.Password != "" {
				err = validate.Var(changes.Password, "gte=4")
				if err != nil {
					response.WriteHeader(500)
					response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
					return
				}
				hash, _ := bcrypt.GenerateFromPassword([]byte(changes.Password), 10)
				author.Password = string(hash)
			}
			authors[index] = author
			json.NewEncoder(response).Encode(authors)
			return
		}
	}
	json.NewEncoder(response).Encode(Author{})
}

func AuthorDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for index, author := range authors {
		if author.Id == params["id"] {
			authors = append(authors[:index], authors[index+1:]...)
			json.NewEncoder(response).Encode(authors)
			return
		}
	}
	json.NewEncoder(response).Encode(Author{})
}
