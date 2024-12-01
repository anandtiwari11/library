package models

import (
	"time"
)

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255); not null" json:"title"`
	ISBN        string    `gorm:"type:varchar(255); unique; not null" json:"isbn"`
	Description string    `gorm:"type:text" json:"description"`
	Authors     []Author  `gorm:"many2many:book_authors;" json:"authors"`
	Available   bool      `json:"available" gorm:"default:true"`
	Users       []User    `json:"-" gorm:"many2many:user_books;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
