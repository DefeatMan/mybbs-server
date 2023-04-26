package models

import "gorm.io/gorm"

// data struct

type Admin struct {
	gorm.Model
	UserId    uint `gorm:"not null;uniqueIndex" json:"userid" binding:"required"`
	AdminPerm uint `gorm:"not null" json:"adminperm" binding:"required"`
}
