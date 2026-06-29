package repository

import (
	"context"
	"errors"
	"fmt"
	"go-functions/internal/response"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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

type VerificationStatus string

type timestamptz string

func (timestamptz) GetGraphQLType() string {
	return "timestamptz"
}

const (
	StatusNoRowExists VerificationStatus = "NO_ROW"
	StatusActiveCode  VerificationStatus = "ACTIVE_CODE_WAIT"
	StatusExpiredCode VerificationStatus = "EXPIRED_CODE_READY"
)

type VerificationAction string

const (
	ActionPasswordReset     VerificationAction = "password_reset"
	ActionEmailVerification VerificationAction = "email_verification"
)

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
			Message:    "Verification Failed.",
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

func (r *HasuraRepository) CheckIfUserExists(ctx context.Context, email string) (bool, error) {
	var query struct {
		Users []struct {
			Id uuid.UUID `graphql:"id"`
		} `graphql:"Users(where: {email: {_eq: $email}}, limit: 1)"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := r.client.Query(ctx, &query, vars); err != nil {
		return false, response.MapDBError(err)
	}

	return len(query.Users) > 0, nil
}

func (r *HasuraRepository) CheckVerificationState(ctx context.Context, email string) (VerificationStatus, *VerificationData, error) {
	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := r.client.Query(ctx, &query, vars); err != nil {
		return StatusNoRowExists, nil, response.MapDBError(err)
	}

	if query.VerificationData.Email == "" {
		return StatusNoRowExists, nil, nil
	}

	if time.Now().After(query.VerificationData.ExpireAt) {
		return StatusExpiredCode, &query.VerificationData, nil
	}

	return StatusActiveCode, &query.VerificationData, nil
}

func (r *HasuraRepository) InsertVerificationRow(ctx context.Context, email, code string, window time.Duration, actionType string) error {
	var mutation struct {
		InsertVerificationDataOne struct {
			Email string `graphql:"email"`
		} `graphql:"insert_VerificationData_one(object: {email: $email, code: $code, expireAt: $expireAt, type: $type})"`
	}

	expireAt := time.Now().Add(window)

	vars := map[string]interface{}{
		"email":    graphql.String(email),
		"code":     graphql.String(code),
		"expireAt": timestamptz(expireAt.Format(time.RFC3339)),
		"type":     graphql.String(actionType),
	}

	if err := r.client.Mutate(ctx, &mutation, vars); err != nil {
		return response.MapDBError(err)
	}

	return nil
}

func (r *HasuraRepository) UpdateUserPassword(ctx context.Context, email, hashedPassword string) error {
	var mutation struct {
		UpdateUsers struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_Users(where: {email: {_eq: $email}}, _set: {password: $password})"`
	}

	vars := map[string]interface{}{
		"email":    graphql.String(email),
		"password": graphql.String(hashedPassword),
	}

	if err := r.client.Mutate(ctx, &mutation, vars); err != nil {
		return response.MapDBError(err)
	}

	if mutation.UpdateUsers.AffectedRows == 0 {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "No active account matches the provided email parameters.",
		}
	}

	return nil
}

func (r *HasuraRepository) ArchiveAndPurgeVerificationRow(ctx context.Context, email, code, actionType, status string) error {

	var transactionMutation struct {
		InsertLogs struct {
			ID string `graphql:"id"`
		} `graphql:"insert_VerificationLogs_one(object: {email: $email, code: $code, type: $type, status: $status})"`

		DeleteData struct {
			Email string `graphql:"email"`
		} `graphql:"delete_VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email":  graphql.String(email),
		"code":   graphql.String(code),
		"type":   graphql.String(actionType),
		"status": graphql.String(status),
	}

	if err := r.client.Mutate(ctx, &transactionMutation, vars); err != nil {
		return response.MapDBError(err)
	}

	return nil
}
