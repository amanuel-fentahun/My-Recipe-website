package main

import (
	"go-functions/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	routes.SetUpRoutes(router)
	router.Run()
}
