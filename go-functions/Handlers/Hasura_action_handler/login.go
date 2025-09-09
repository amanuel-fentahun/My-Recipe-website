package hasuraactionhandler

import (
	"context"
	. "go-functions/Utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type LoginActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            map[string]interface{} `json:"input"`
}

func LoginHandler(c *gin.Context) {
	var actionPayload LoginActionPayload

	if err := c.ShouldBindJSON(&actionPayload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	input := actionPayload.Input

	credintails, ok := input["credentails"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid credentials format"})
		return
	}

	email, ok := credintails["email"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Email is required"})
	}

	password, ok := credintails["password"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Password is required"})
	}

	var query struct {
		UsersAggrigation struct {
			Nodes []struct {
				Email     string    `graphql:"email"`
				Password  string    `graphql:"password"`
				Name      string    `graphql:"name"`
				ID        uuid.UUID `graphql:"id"`
				AvaterURL string    `graphql:"avater_url"`
			} `graphql:"nodes"`
		} `graphql:"Users_aggregate(where: {email: {_eq: $email}})"`
	}

	variables := map[string]interface{}{
		"email": graphql.String(email),
	}

	client := graphql.NewClient(os.Getenv("HASURA_GRAPHQL_ENDPOINT"), nil)
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	if len(query.UsersAggrigation.Nodes) == 0 {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}
	user := query.UsersAggrigation.Nodes[0]

	match := CheckPasswordHash(password, user.Password)

	if !match {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := CreateToken(user.Name, user.Email, user.ID, user.AvaterURL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(200, gin.H{
		"accessToken": token,
	})

}
