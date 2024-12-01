package routes

import (
	"github.com/anandtiwari11/library-management/controllers"
	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	router.GET("/get-books", controllers.GetBooks)
	router.GET("/get-books/:id", controllers.GetBook)
	router.POST("/new-books", controllers.PostBook)
}