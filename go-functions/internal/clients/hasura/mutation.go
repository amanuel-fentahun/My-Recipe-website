package hasura

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

func MarkEmailVerfied(email string) error {
	var Mutation struct {
		UpdateUsers struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_Users(where: {email: {_eq: $email}}, _set: {isVerified: true})"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := NewGraphqlClient().Mutate(context.Background(), &Mutation, vars); err != nil {
		return err
	}

	if Mutation.UpdateUsers.AffectedRows == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}

	return nil
}
