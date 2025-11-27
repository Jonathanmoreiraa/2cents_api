package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Expense struct {
	ID          int             `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int             `json:"user_id,omitempty" gorm:"not null"`
	User        User            `json:"-" gorm:"constraint:OnUpdate:CASCADE"`
	CategoryID  int             `json:"category_id,omitempty" gorm:"not null"`
	Category    Category        `json:"-" gorm:"constraint:OnUpdate:CASCADE"`
	SavingID    *int            `json:"saving_id,omitempty" gorm:"index;constraint:OnDelete:CASCADE"`
	Saving      *Saving         `json:"-" gorm:"constraint:OnUpdate:CASCADE"`
	Description string          `json:"description" gorm:"not null;type:varchar(255)"`
	Value       decimal.Decimal `json:"value" gorm:"not null;type:decimal(19,2)"`
	DueDate     *time.Time      `json:"due_date" gorm:"type:date"`
	Paid        int             `json:"paid" gorm:"not null;type:tinyint(1);default:0"`
	CreatedAt   time.Time       `json:"created" gorm:"not null;"`
	UpdatedAt   time.Time       `json:"modified" gorm:"not null;"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}
