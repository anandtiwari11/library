package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anandtiwari11/library-management/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func RequireAuth(c *gin.Context) {
    tokenString, err := c.Cookie("Authorization") 
    if err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("Anand-Tiwari"), nil
    })

    if err != nil || !token.Valid {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        fmt.Println("Token Claims: ", claims)

        db := c.MustGet("library").(*gorm.DB)
        var user models.User
        if err := db.Preload("BorrowedBooks").First(&user, "email = ?", claims["email"]).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
            }
            return
        }        

        if user.ID == 0 {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        c.Set("user", user)
        c.Next()
    } else {
        c.AbortWithStatus(http.StatusUnauthorized)
    }
}