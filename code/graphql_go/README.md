![](../../golangnetwork-logo.png)

# GraphQL em Go

[**Cleuton Sampaio**](https://linkedin.com/in/cleutonsampaio) - Me siga!
Repositório: [https://github.com/cleuton/golang-network](https://github.com/cleuton/golang-network).

## GraphQL é baseado em grafos

Você pode solicitar qualquer pedaço do grafo de informação.
Usamos requests GET ou POST para enviar queries ao servidor GraphQL. É mais comum e seguro
enviarmos utilizando POST, que não é idempotente.

## Esquema

Leia a documentação oficial do GraphQL. A implementação em Golang pode ser um pouco diferente, então vamos nos ater ao site oficial.

Imagine um esquema como esse:
```graphql
type Query {
    url(name: String!) : String
}
```

Esta é a linguagem **SDL** - Schema Definition Language do GraphQL. Estamos declarando um tipo de dados "Query" que tem um campo "url". Os campos podem ou não ter argumentos. Neste caso, para obter a URL, precisamos passar o NOME do site. O nome não pode ser null (note a exclamação após o tipo dedados).

Se quisermos saber a URL de um site, fazemos uma query assim:

```sdl
{
    url(name: "github")
}
``` 

E podemos enviar esta query utilizando POST ou GET. No caso de GET, poderíamos enviar algo assim:

```shell
http://myapi/graphql?query={url(name:"github")}
```

Em um Request POST, podemos enviar um content-type application/json no body com este formato:

```sdl
{
    "query" : "{url(name: \"github\")}"
}
```

O resultado poderia ser algo assim:

```json
{
    "data" : {"url":"http://github.com"}
}
```

## Mutations

Se quisermos alterar algo no servidor, então precisamos definir uma mutação (mutation). Por exemplo:

```sdl
type Mutation {
    addSite(name: String!, url: String!): Boolean!
}
```

A mutação "addSite" recebe dois argumentos e retorna um Boolean (True / False), indicando se foi possível ou não adicionar o site.
Para executar, basta enviar isso em um POST:

```json
{
    "query" : "mutation {addSite(\"netflix.com\")}"
}
```

A resposta pode ser algo assim:

```json
{
    "data" : {"addSite":true}
}
``` 

No repositório você pode ver o código em **Go** dessa implementação: 

```go
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
```



