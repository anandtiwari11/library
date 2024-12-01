package routes

import (
	"github.com/anandtiwari11/library-management/controllers"
	"github.com/gin-gonic/gin"
)


func AuthorRoutes(router *gin.Engine) {
	router.GET("/authors", controllers.GetAllAuthors)
	router.GET("/author/:id/book", controllers.AuthorBooks)
	router.POST("authors/new", controllers.NewAuthor)
}