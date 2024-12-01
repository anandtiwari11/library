package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anandtiwari11/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var user models.User

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
		return
	}

	if err := db.Preload("BorrowedBooks").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the User"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var request struct {
		Name            string        `json:"name"`
		Email           string        `json:"email" gorm:"unique"`
		SubscriptionEnd time.Time     `json:"subscription_end"`
		BorrowedBooks   []models.Book `json:"borrowed_books" gorm:"many2many:user_books;"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Data"})
		return
	}
	user := models.User{
		Name:            request.Name,
		Email:           request.Email,
		SubscriptionEnd: request.SubscriptionEnd,
		BorrowedBooks:   request.BorrowedBooks,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func Borrow(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var user models.User
	var book models.Book

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID"})
		return
	}
	
	if err := db.First(&book, bookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the book"})
		}
		return
	}

	if err := db.Preload("BorrowedBooks").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the user"})
		}
		return
	}

	if !book.Available {
		c.JSON(http.StatusConflict, gin.H{"error": "Book not available"})
		return
	}


	book.Available = false
	user.BorrowedBooks = append(user.BorrowedBooks, book)


	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to borrow book"})
		return
	}

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book borrowed successfully", "user": user})
}
