package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Email       string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	PhoneNumber string         `gorm:"type:varchar(20)" json:"phone_number"`
	Position    string         `gorm:"type:varchar(100)" json:"position"`
	Department  string         `gorm:"type:varchar(100)" json:"department"`
	Salary      float64        `gorm:"type:decimal(10,2)" json:"salary"`
	HireDate    time.Time      `gorm:"not null" json:"hire_date"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Employee) TableName() string {
	return "employees"
}
