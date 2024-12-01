package controllers

import (
	"net/http"
	"strconv"

	"github.com/anandtiwari11/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBooks(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var books []models.Book

	if err := db.Preload("Authors").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var book models.Book

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := db.Preload("Authors").First(&book, bookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

func PostBook(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var request struct {
		Title       string        `json:"title"`
		ISBN        string        `json:"isbn"`
		Description string        `json:"description"`
		Available   bool          `json:"available" gorm:"default:true"`
		Users       []models.User `json:"-" gorm:"many2many:user_books;"`
		AuthorIDs   []uint        `json:"author_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var authors []models.Author
	if err := db.Where("id IN ?", request.AuthorIDs).Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}

	book := models.Book{
		Title:       request.Title,
		ISBN:        request.ISBN,
		Description: request.Description,
		Authors:     authors,
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}
