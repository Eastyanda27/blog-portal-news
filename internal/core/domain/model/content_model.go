package model

import "time"

type Content struct {
	ID          int64      `gorm:"id"`
	CategoryID  int64      `gorm:"category_id"`
	CreatedByID int64      `gorm:"created_by_id"`
	Title       string     `gorm:"title"`
	Excerpt     string     `gorm:"excerpt"`
	Description string     `gorm:"description"`
	Tags        string     `gorm:"tags"`
	Status      string     `gorm:"status"`
	Category    Category   `gorm:"foreignKey:CategoryID"`
	User        User       `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
}
