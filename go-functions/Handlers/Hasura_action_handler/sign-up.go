package hasuraactionhandler

import (
	"context"
	"fmt"
	"os"

	. "go-functions/Utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type SignupActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            map[string]interface{} `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

func SignUpHandler(c *gin.Context) {

	var actionPayload SignupActionPayload

	if err := c.ShouldBindJSON(&actionPayload); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// creating a new request client using go-graphql-client
	client := graphql.NewClient(os.Getenv("HASURA_GRAPHQL_ENDPOINT"), nil)

	inputs := actionPayload.Input

	credintails, ok := inputs["inputs"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid credentials format"})
		return
	}

	name, ok1 := credintails["fullName"].(string)
	email, ok2 := credintails["email"].(string)
	password, ok3 := credintails["password"].(string)
	avater_url, ok := credintails["avater_url"].(string)

	// try to check if all required fields are present
	if !ok1 || !ok2 || !ok3 || name == "" || email == "" || password == "" {
		c.JSON(400, gin.H{
			"error": "fullName, email and password are required",
		})
		return
	}

	// GraphQL query for select one user by email
	var query struct {
		UsersAggrigation struct {
			Nodes []struct {
				ID uuid.UUID `graphql:"id"`
			} `graphql:"nodes"`
		} `graphql:"Users_aggregate(where: {email: {_eq: $email}})"`
	}

	variables := map[string]interface{}{
		"email": graphql.String(email),
	}

	// try to fetch user with the given email
	if err := client.Query(context.Background(), &query, variables); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to check existing user",
		})
		return
	}

	// check if the user already exists
	if len(query.UsersAggrigation.Nodes) != 0 {
		c.JSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	// hash the user password
	hashed_password, err := HashPassword(password)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//Send A verification Email
	code := GenerateRandomString(6)

	subject := "Verify your Email"
	body := "<p> You need to verfiy your email address to continue using your <strong>Tafach Kitchen</strong> account.</p>"
	body += "<p>Enter the following code to verify your email address:</p>"
	body += fmt.Sprintf("<h3>%s</h3>", code)
	body += "<p>If you did not create this account, please ignore this email.</p>"
	body += "<p>Thanks,<br/>The Tafach Kitchen Team</p>"

	if err := SendEmail(email, subject, body); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to send verification email",
		})
		return
	}

	// Graphql mutation for signup
	var SignupMutation struct {
		InsertUsersOne struct {
			ID uuid.UUID `graphql:"id"`
		} `graphql:"insert_Users_one(object: {name: $name, email: $email, password: $password, avater_url: $avater_url})"`
	}

	varaibles2 := map[string]interface{}{
		"name":       graphql.String(name),
		"email":      graphql.String(email),
		"password":   graphql.String(hashed_password),
		"avater_url": graphql.String(avater_url),
	}

	// try to create a new user
	if err := client.Mutate(context.Background(), &SignupMutation, varaibles2); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully, Please verify your email",
	})

}
