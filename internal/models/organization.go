package models

import "time"

type Organization struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"unique;not null"`
	Email     string     `gorm:"unique;not null"`
	Contracts []Contract `gorm:"foreignKey:OrganizationID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Organization) TableName() string {
	return "organizations"
}
