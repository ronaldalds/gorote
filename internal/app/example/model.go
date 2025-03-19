package example

import (
	"gorm.io/gorm"
)

type Example struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;size:100;not null" validate:"required,min=3,max=100"`
}