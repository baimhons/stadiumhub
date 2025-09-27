package team

import "time"

type Team struct {
	ID             int       `gorm:"primaryKey"`
	Name           string    `gorm:"not null,unique"`
	ShortName      string    `gorm:"not null,unique"`
	TLA            string    `gorm:"not null"`
	Address        string    `gorm:"not null"`
	Venue          string    `gorm:"not null,unique"`
	ViewerCapacity int       `gorm:"not null"`
	Price          float32   `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"not null default:current_timestamp"`
}
