package config

import (
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
)

func NewGraphqlClient() *graphql.Client {

	client := graphql.NewClient(os.Getenv("HASURA_GRAPHQL_ENDPOINT"), nil).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))
		})

	return client
}
