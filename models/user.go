package models

import "time"

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name"`
	Email           string     `json:"email" gorm:"unique"`
	SubscriptionEnd time.Time  `json:"subscription_end"`
	BorrowedBooks   []Book     `json:"borrowed_books" gorm:"many2many:user_books;"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
