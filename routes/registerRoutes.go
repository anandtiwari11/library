package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	BookRoutes(router)
	AuthorRoutes(router)
	UserRoutes(router)
}