package models

import (
	"time"
)

// data struct

type CommentModel struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:idx_postid_createtime,priority:10,sort:asc; index:idx_postid_agreenum_createtime,priority:10,sort:asc"`
	UpdatedAt time.Time
}

type Comment struct {
	//gorm.Model
	CommentModel
	PostId   uint   `gorm:"not null;index:idx_postid_createtime,priority:8;index:id_postid_agreenum_createtime,priority:8" json:"postid" binding:"required"`
	LinkId   uint   `gorm:"not null" json:"linkid" binding:"required"`
	UserId   uint   `gorm:"not null" json:"userid"`
	Content  string `gorm:"type:text;not null" json:"content" binding:"required"`
	AgreeNum uint   `gorm:"not null; index:idx_postid_agreenum_createtime,priority:9,sort:desc"`
	StarNum  uint   `gorm:"not null"`
}

// post form

type CommentAppendForm struct {
	CommentId uint   `json:"commentid" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
