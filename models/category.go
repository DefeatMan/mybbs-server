package models

import "gorm.io/gorm"

// data struct

type Category struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null;unique;uniqueIndex" json:"name" binding:"required"`
	FollowNum uint   `gorm:"not null;index:,sort:desc"`
}

// post form

type CategoryCreateForm struct {
	Name string `json:"name" binding:"required"`
}

type CategoryRenameForm struct {
	CatagoryId uint   `json:"categoryid" binding:"required"`
	Name       string `json:"name" binding:"required"`
}
