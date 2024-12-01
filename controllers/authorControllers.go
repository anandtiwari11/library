package controllers

import (
	"net/http"
	"strconv"

	"github.com/anandtiwari11/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllAuthors(c *gin.Context) {
	var authors []models.Author
	db, exists := c.Get("library")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	database, ok := db.(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database instance"})
		return
	}
	if err := database.Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}
	c.JSON(http.StatusOK, authors)
}

func AuthorBooks(c *gin.Context) {
	var author models.Author
	db := c.MustGet("library").(*gorm.DB)

	authorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	if err := db.Preload("Books.Authors").First(&author, authorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch author data"})
		}
		return
	}

	c.JSON(http.StatusOK, author.Books)
}

func NewAuthor(c *gin.Context) {
	var author models.Author
	db := c.MustGet("library").(*gorm.DB)
	
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}
	
	if err := db.Create(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, author)
}