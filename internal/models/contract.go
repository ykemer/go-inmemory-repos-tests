package models

import "time"

type Contract struct {
	ID             uint         `gorm:"primaryKey"`
	OrganizationID uint         `gorm:"not null"`
	Title          string       `gorm:"not null"`
	Description    string       `gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	StartDate      time.Time
	EndDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Contract) TableName() string {
	return "contracts"
}
