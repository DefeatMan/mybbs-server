package models

import (
	"time"
)

// data struct

type PostModel struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:idx_catagory_createtime,priority:10,sort:desc; index:idx_category_starnum_createtime,priority:10,sort:desc"`
	UpdatedAt time.Time
}

type Post struct {
	PostModel
	CategoryId uint   `gorm:"not null;index:idx_catagory_createtime,priority:8;index:idx_category_starnum_createtime,priority:8"`
	Title      string `gorm:"type:varchar(100);not null"`
	UserId     uint   `gorm:"not null;index"`
	CommentId  uint   `gorm:"not null"`
	StarNum    uint   `gorm:"not null;index:,desc;index:idx_category_starnum_createtime,priority:9,sort:desc"`
	LockFlag   uint   `gorm:"not null"`
}

// post form
type PostCreateForm struct {
	CategoryId uint   `json:"categoryid" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// flag

const BaseTextFlag uint = 0
