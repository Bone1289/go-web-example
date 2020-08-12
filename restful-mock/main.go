package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

var articles = []Article{
	{
		Id:      "article-1",
		Author:  "author-1",
		Title:   "This is an Example Article",
		Content: "This is some sample content",
	},
}

var authors = []Author{
	{
		Id:        "author-1",
		Firstname: "Nicolas",
		Lastname:  "Raboy",
		Username:  "nraboy",
		Password:  "$2a$10$0OtFx9DSi5x.bnjx28f4Xu1pkURjYVnTvgFnvoxIdyXambjSyLQhW",
	},
	{
		Id:        "author-2",
		Firstname: "Maria",
		Lastname:  "Raboy",
		Username:  "mraboy",
		Password:  "$2a$10$0OtFx9DSi5x.bnjx28f4Xu1pkURjYVnTvgFnvoxIdyXambjSyLQhW",
	},
}

func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Write([]byte(`{ "message": "Hello World" }`))
}

func main() {
	fmt.Println("Starting application...")
	router := mux.NewRouter()
	router.HandleFunc("/", RootEndpoint).Methods("GET")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/author", RegisterEndpoint).Methods("POST")
	router.HandleFunc("/authors", AuthorRetrieveAllEndpoint).Methods("GET")
	router.HandleFunc("/author/{id}", AuthorRetrieveEndpoint).Methods("GET")
	router.HandleFunc("/author/{id}", AuthorUpdateEndpoint).Methods("PUT")
	router.HandleFunc("/author/{id}", AuthorDeleteEndpoint).Methods("DELETE")
	router.HandleFunc("/article", ValidateMiddleware(ArticleCreateEndpoint)).Methods("POST")
	router.HandleFunc("/articles", ArticleRetrieveAllEndpoint).Methods("GET")
	router.HandleFunc("/article/{id}", ArticleRetrieveEndpoint).Methods("GET")
	router.HandleFunc("/article/{id}", ValidateMiddleware(ArticleUpdateEndpoint)).Methods("PUT")
	router.HandleFunc("/article/{id}", ValidateMiddleware(ArticleDeleteEndpoint)).Methods("DELETE")
	headers := handlers.AllowedHeaders(
		[]string{
			"X-Requested-With",
			"Content-Type",
			"Authorization",
		},
	)
	methods := handlers.AllowedMethods(
		[]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
	)
	origins := handlers.AllowedOrigins(
		[]string{
			"*",
		},
	)
	http.ListenAndServe(
		":12345",
		handlers.CORS(headers, methods, origins)(router),
	)
}
