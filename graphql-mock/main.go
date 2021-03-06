package main

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
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

var rootMutation *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createArticle": &graphql.Field{
			Type: graphql.NewList(articleType),
			Args: graphql.FieldConfigArgument{
				"article": &graphql.ArgumentConfig{
					Type: articleInputType,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var article Article
				mapstructure.Decode(params.Args["article"], &article)

				decoded, err := ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(article)
				if err != nil {
					return nil, err
				}

				article.Id = uuid.Must(uuid.NewV4()).String()
				article.Author = decoded.(CustomJWTClaims).Id
				articles = append(articles, article)
				return articles, nil
			},
		},
		"updateAuthor": &graphql.Field{
			Type: graphql.NewList(authorType),
			Args: graphql.FieldConfigArgument{
				"author": &graphql.ArgumentConfig{
					Type: authorInputType,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var changes Author
				mapstructure.Decode(params.Args["author"], &changes)
				validate := validator.New()
				err := validate.StructExcept(changes, "Firstname", "Lastname", "Username", "Password", "Id")
				if err != nil {
					return nil, err
				}

				for index, author := range authors {
					if author.Id == changes.Id {
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
								return nil, err
							}
							hash, _ := bcrypt.GenerateFromPassword([]byte(changes.Password), 10)
							author.Password = string(hash)
						}
						authors[index] = author
						return authors, nil
					}
				}
				return nil, nil
			},
		},
		"deleteAuthor": &graphql.Field{
			Type: graphql.NewList(articleType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				for index, author := range authors {
					if author.Id == id {
						authors = append(authors[:index], authors[index+1:]...)
						return authors, nil
					}
				}
				return nil, nil
			},
		},
	},
})

type CustomJWTClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

var JwtSecret []byte = []byte("testtest")

type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	router := mux.NewRouter()
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	router.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		var payload GraphQLPayload
		json.NewDecoder(request.Body).Decode(&payload)
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.Query,
			VariableValues: payload.Variables,
			Context:        context.WithValue(context.Background(), "token", request.URL.Query().Get("token")),
		})
		json.NewEncoder(response).Encode(result)
	})
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/author", RegisterEndpoint).Methods("POST")

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
