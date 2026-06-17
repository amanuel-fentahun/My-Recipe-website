package repository

import (
	"context"
	"errors"
	"fmt"
	"go-functions/internal/response"
	"net/http"
	"os"
	"time"

	"github.com/hasura/go-graphql-client"
)

type HasuraRepository struct {
	client *graphql.Client
}

func NewHasuraRepository() *HasuraRepository {
	endpoint := os.Getenv("HASURA_GRAPHQL_ENDPOINT")
	adminSecret := os.Getenv("HASURA_ADMIN_SECRET")

	client := graphql.NewClient(endpoint, nil).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-hasura-admin-secret", adminSecret)
		})

	return &HasuraRepository{client: client}
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

func (r *HasuraRepository) FetchVerificationDataByEmail(ctx context.Context, email string) (*VerificationData, error) {
	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := r.client.Query(ctx, &query, vars); err != nil {
		return nil, response.MapDBError(err)
	}

	if query.VerificationData.Email == "" {
		err := errors.New("empty data result block returned")
		return nil, &response.AppError{
			HTTPStatus: http.StatusNotFound,
			Code:       response.CodeInvalidInput,
			Message:    "ErrUserNotFound",
			RawError:   err,
		}
	}

	return &query.VerificationData, nil
}

func (r *HasuraRepository) MarkEmailVerified(ctx context.Context, email string) error {
	var mutation struct {
		UpdateUsers struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_Users(where: {email: {_eq: $email}}, _set: {isVerified: true})"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := r.client.Mutate(ctx, &mutation, vars); err != nil {
		return response.MapDBError(err)
	}

	if mutation.UpdateUsers.AffectedRows == 0 {
		err := fmt.Errorf("no user found with email %s", email)
		return &response.AppError{
			HTTPStatus: http.StatusNotFound,
			Code:       response.CodeInvalidInput,
			Message:    "Failed to update verification status: profile not found",
			RawError:   err,
		}
	}

	return nil
}
