package model

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	Id            uint    `gorm:"primarykey"`
	TransactionId string  `gorm:"unique;not null"`
	UserId        string  `gorm:"not null;index"`
	Amount        float64 `gorm:"not null"`
	Type          string  `gorm:"not null"` // e.g., "credit", "debit"
	Status        string  `gorm:"not null"` // e.g., "pending", "completed", "failed"
	Description   string  `gorm:""`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	
	// Relationship with User
	User User `gorm:"foreignKey:UserId;references:UserId"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}