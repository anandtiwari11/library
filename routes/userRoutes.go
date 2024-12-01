package routes

import (
	"github.com/anandtiwari11/library-management/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/user/:id", controllers.GetUser)
	router.POST("/create-user", controllers.CreateUser)
	router.POST("/borrow/:userID/:bookID", controllers.Borrow)
	router.PUT("/return/:userID/:bookID", controllers.ReturnBook)
}