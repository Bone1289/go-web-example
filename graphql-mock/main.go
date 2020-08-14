package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"net/http"
)

var authors []Author = []Author{
	Author{
		Id:        "author-1",
		Firstname: "Nicolas",
		Lastname:  "Raboy",
		Username:  "nraboy",
		Password:  "$2a$10$0OtFx9DSi5x.bnjx28f4Xu1pkURjYVnTvgFnvoxIdyXambjSyLQhW",
	},
	Author{
		Id:        "author-2",
		Firstname: "Maria",
		Lastname:  "Raboy",
		Username:  "mraboy",
		Password:  "$2a$10$0OtFx9DSi5x.bnjx28f4Xu1pkURjYVnTvgFnvoxIdyXambjSyLQhW",
	},
}

var articles []Article = []Article{
	Article{
		Id:      "article-1",
		Author:  "author-1",
		Title:   "This is an Example Article",
		Content: "This is some sample content",
	},
}

var rootQuery *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"authors": &graphql.Field{
			Type: graphql.NewList(authorType),
			Resolve: func(param graphql.ResolveParams) (interface{}, error) {
				return authors, nil
			},
		},
		"author": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(param graphql.ResolveParams) (interface{}, error) {
				id := param.Args["id"].(string)
				for _, author := range authors {
					if author.Id == id {
						return author, nil
					}
				}
				return nil, nil
			},
		},
		"articles": &graphql.Field{
			Type: graphql.NewList(articleType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return articles, nil
			},
		},
		"article": &graphql.Field{
			Type: articleType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(param graphql.ResolveParams) (interface{}, error) {
				id := param.Args["id"].(string)
				for _, article := range articles {
					if article.Id == id {
						return article, nil
					}
				}
				return nil, nil
			},
		},
	},
})

type GraphQLPayload struct {
	Query string `json:"query"`
}

func main() {
	router := mux.NewRouter()
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: nil,
	})

	router.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		var payload GraphQLPayload
		json.NewDecoder(request.Body).Decode(&payload)
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: payload.Query,
		})
		json.NewEncoder(response).Encode(result)
	})
	http.ListenAndServe(":12345", router)
}
