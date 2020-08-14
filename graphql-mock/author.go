package main

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Author struct {
	Id        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

var authorType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"firstname": &graphql.Field{
			Type: graphql.String,
		},
		"lastname": &graphql.Field{
			Type: graphql.String,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authorInputType *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AuthorInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"firstname": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"lastname": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"username": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"password": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

func RegisterEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var author Author
	json.NewDecoder(request.Body).Decode(&author)
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
	for _, author := range authors {
		if author.Username == data.Username {
			err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(data.Password))
			if err != nil {
				response.WriteHeader(500)
				response.Write([]byte(`{ "message": "invalid password" }`))
				return
			}
			json.NewEncoder(response).Encode(author)
			return
		}
	}
	response.Write([]byte(`{ "message": "invalid username" }`))
}
