package controllers

import (
	"net/http"

	"github.com/anandtiwari11/library-management/helpers"
	"github.com/anandtiwari11/library-management/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var input, user models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := db.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"Error": "User Already Exist"})
		return
	}
	hashedString, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	user = models.User{
		Name :			 input.Name,
		Email:           input.Email,
		Password:        string(hashedString),
		SubscriptionEnd: user.SubscriptionEnd,
		BorrowedBooks:   user.BorrowedBooks,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, "User Created Succesfully")
}

func Login(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"Error": "User Not Exist"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed ot compare passwords")
		return
	}
	tokenString, err := helpers.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to Generate Token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 72, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message" : "Login Successfull",
	})
}

func GetUser(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"user": user})
}

func Logout(c *gin.Context) {
    c.SetCookie("Authorization", "", -1, "", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
