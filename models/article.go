package models

import (
	"time"
)

// Article merepresentasikan model artikel/berita dalam sistem
type Article struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	AuthorID  uint      `gorm:"not null" json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
