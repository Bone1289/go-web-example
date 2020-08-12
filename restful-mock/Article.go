package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	uuid "github.com/satori/go.uuid"
	validator "gopkg.in/go-playground/validator.v9"

	"net/http"
)

type Article struct {
	Id      string `json:"id,omitempty" validate:"omitempty,uuid"`
	Author  string `json:"author,omitempty" validate:"omitempty"`
	Title   string `json:"title,omitempty" validate:"required"`
	Content string `json:"content,omitempty" validate:"required"`
}

func ArticleCreateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var article Article
	token := context.Get(request, "decoded").(CustomJWTClaims)
	json.NewDecoder(request.Body).Decode(&article)

	validate := validator.New()

	err := validate.Struct(article)

	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{"message":"` + err.Error() + `" }"`))
		return
	}

	article.Id = uuid.Must(uuid.NewV4()).String()
	article.Author = token.Id
	articles = append(articles, article)
	json.NewEncoder(response).Encode(article)
}

func ArticleRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(articles)
}

func ArticleRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for _, article := range articles {
		if article.Id == params["id"] {
			json.NewEncoder(response).Encode(article)
			return
		}
	}
	json.NewEncoder(response).Encode(Article{})
}

func ArticleUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var changes Article
	params := mux.Vars(request)
	token := context.Get(request, "decoded").(CustomJWTClaims)
	json.NewDecoder(request.Body).Decode(&changes)
	validate := validator.New()
	err := validate.StructExcept(changes, "Title", "Content")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for index, article := range articles {
		if article.Id == params["id"] && article.Author == token.Id {
			if changes.Title != "" {
				article.Title = changes.Title
			}
			if changes.Content != "" {
				article.Content = changes.Content
			}
			articles[index] = article
			json.NewEncoder(response).Encode(articles)
			return
		}
	}
	json.NewEncoder(response).Encode(Article{})
}

func ArticleDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	token := context.Get(request, "decoded").(CustomJWTClaims)
	for index, article := range articles {
		if article.Id == params["id"] && article.Author == token.Id {
			articles = append(articles[:index], articles[index+1:]...)
			json.NewEncoder(response).Encode(articles)
			return
		}
	}
	json.NewEncoder(response).Encode(Article{})
}
