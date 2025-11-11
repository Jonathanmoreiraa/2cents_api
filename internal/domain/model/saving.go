package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TODO: ver para adicionar um valor inicial, porque o usuário pode já ter um valor guardado
type Saving struct {
	ID              int             `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int             `json:"user_id,omitempty" gorm:"not null"`
	User            User            `json:"-" gorm:"constraint:OnUpdate:CASCADE"`
	Description     string          `json:"description" gorm:"not null;type:varchar(255)"`
	Goal            decimal.Decimal `json:"goal" gorm:"not null;type:decimal(19,2)"`
	Accumulated     decimal.Decimal `json:"accumulated" gorm:"not null;type:decimal(19,2)"`
	IsEmergencyFund int             `json:"is_emergency_fund" gorm:"type:tinyint(1);not null;default:0"`
	Priority        int             `json:"priority" gorm:"not null"`
	CreatedAt       time.Time       `json:"created" gorm:"not null;"`
	UpdatedAt       time.Time       `json:"modified" gorm:"not null;"`
	DeletedAt       gorm.DeletedAt  `gorm:"index"`
}
