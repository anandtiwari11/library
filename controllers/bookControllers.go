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

	if len(authors) != len(request.AuthorIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more authors not found"})
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

	for i := range authors {
		authors[i].Books = append(authors[i].Books, book)
		if err := db.Save(&authors[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author with book"})
			return
		}
	}

	c.JSON(http.StatusCreated, book)
}

func Borrow(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var user models.User
	var book models.Book

	userID, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "unathorized Access"})
	}

	bookID, err := strconv.Atoi(c.Param("bookID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID"})
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

func ReturnBook(c *gin.Context) {
	db := c.MustGet("library").(*gorm.DB)
	var user models.User
	var book models.Book

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unathorized User"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID"})
		return
	}
	if err := db.Preload("BorrowedBooks").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	bookFound := false
	for i, borrowedBook := range user.BorrowedBooks {
		if borrowedBook.ID == book.ID {
			user.BorrowedBooks = append(user.BorrowedBooks[:i], user.BorrowedBooks[i+1:]...)
			bookFound = true
			break
		}
	}

	if !bookFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found in user's borrowed books"})
		return
	}

	book.Available = true

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if err := db.Model(&user).Association("BorrowedBooks").Delete(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user-book association"})
		return
	}

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully", "user": user, "book": book})
}