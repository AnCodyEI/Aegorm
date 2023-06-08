package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	GuardName   string `gorm:"size:255;not null;index"`
	Description string `gorm:"size:255;"`
}

type Roles []Role

func (r *Role) BeforeCreate(*gorm.DB) (err error) {
	r.GuardName = slug.Make(r.Name)
	return
}
