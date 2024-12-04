package routes

import (
	"github.com/anandtiwari11/library-management/controllers"
	"github.com/anandtiwari11/library-management/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.GetUser)
	router.POST("/logout", middleware.RequireAuth, controllers.Logout)
}