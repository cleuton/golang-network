package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

// Schema creation

type Url struct {
	Name    string `json:"name"`
	SiteUrl string `json:"siteurl"`
}

var (
	UrlList []Url

	// The GraphQL "Url" type:
	UrlType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Url",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"siteurl": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// The main GraphQL query:
	RootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{

			/*
			   {"query" : "{url(name:\"google\") {name siteurl}}"}
			*/
			"url": &graphql.Field{
				Type:        UrlType,
				Description: "Get single URL",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetResolver,
			},

			/*
			   {"query": "{urllist {name siteurl}}"}
			*/
			"urllist": &graphql.Field{
				Type:        graphql.NewList(UrlType),
				Description: "List of urls",
				Resolve:     GetListResolver,
			},
		},
	})

	// The mutation to add or delete Url:
	RootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			/*
				{ "query": "mutation { createUrl(name:\"youtube\", siteurl:\"youtube.com\") { name siteurl } }" }
			*/
			"createUrl": &graphql.Field{
				Type:        UrlType, // the return type for this field
				Description: "Create new Url",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"siteurl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: CreateNewUrlResolve,
			},
			/*
				{ "query": "mutation { deleteUrl(name:\"youtube\") { name siteurl } }" }
			*/
			"deleteUrl": &graphql.Field{
				Type:        UrlType, // the return type for this field
				Description: "Delete an Url",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: DeleteUrlResolve,
			},
		},
	})

	// The GraphQL schema:
	UrlSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    RootQuery,
		Mutation: RootMutation,
	})
)

func init() {
	UrlList = append(UrlList, Url{Name: "google", SiteUrl: "http://google.com"})
}

// GraphQL resolvers:

var GetResolver = func(params graphql.ResolveParams) (interface{}, error) {

	idQuery, isOK := params.Args["name"].(string)
	if isOK {
		// Search for el with id
		for _, url := range UrlList {
			if url.Name == idQuery {
				return url, nil
			}
		}
	}

	return Url{}, nil
}

var GetListResolver = func(p graphql.ResolveParams) (interface{}, error) {
	return UrlList, nil
}

var CreateNewUrlResolve = func(params graphql.ResolveParams) (interface{}, error) {

	nameProp, _ := params.Args["name"].(string)
	siteurlProp, _ := params.Args["siteurl"].(string)

	newUrl := Url{Name: nameProp, SiteUrl: siteurlProp}
	UrlList = append(UrlList, newUrl)

	return newUrl, nil
}

var DeleteUrlResolve = func(params graphql.ResolveParams) (interface{}, error) {

	nameUrl, _ := params.Args["name"].(string)
	urlDeletada := Url{}

	for i := 0; i < len(UrlList); i++ {
		if UrlList[i].Name == nameUrl {
			urlDeletada = UrlList[i]
			UrlList = append(UrlList[:i], UrlList[i+1:]...)
			break
		}
	}
	// Return affected todo
	return urlDeletada, nil
}

// Http handle functions:

type PostData struct {
	Query string `json:"query"`
}

func ProcessGraphQL(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var query PostData
	if err := decoder.Decode(&query); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}
	result := graphql.Do(graphql.Params{
		Context:       r.Context(),
		Schema:        UrlSchema,
		RequestString: query.Query,
	})

	WriteResponse(http.StatusOK, result, w)

}

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/graphql", ProcessGraphQL).Methods("POST")
	err := http.ListenAndServe("localhost:8080", router)
	fmt.Println(err)
}
