package routes

import (
	"github.com/anandtiwari11/library-management/controllers"
	"github.com/anandtiwari11/library-management/middleware"
	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	router.GET("/books", controllers.GetBooks)
	router.GET("/book/:id", controllers.GetBook)
	router.POST("/addbook", controllers.PostBook)
	router.POST("/borrow/:bookID", middleware.RequireAuth, controllers.Borrow)
	router.PUT("/return/:bookID", middleware.RequireAuth, controllers.ReturnBook)
}