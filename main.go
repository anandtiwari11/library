package main

import (
	"log"

	"github.com/anandtiwari11/library-management/helpers"
	"github.com/anandtiwari11/library-management/routes"
	"github.com/gin-gonic/gin"
)


func main() {
	helpers.ConnectDB()
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("library", helpers.DB)
		c.Next()
	})
	routes.RegisterRoutes(router)
	log.Fatal(router.Run(":8080"))
}