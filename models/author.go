package models

import "time"

type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Books     []Book    `gorm:"many2many:book_authors" json:"books"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
