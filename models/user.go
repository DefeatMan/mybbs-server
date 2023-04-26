package models

import (
	"gorm.io/gorm"
	"time"
)

// data struct

type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(20);not null;unique;uniqueIndex"`
	Email         string `gorm:"size:255;not null;unique;uniqueIndex"`
	Password      string `gorm:"size:255;not null"`
	AgreeNum      uint   `gorm:"not null"`
	StarNum       uint   `gorm:"not null"`
	LastLoginIpv4 string `gorm:"type:varchar(20)"`
}

// follow & star

type UserStarPost struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:,desc"`
    UserId    uint      `gorm:"not null;uniqueIndex:idx_userid_star_postid,priority:8" json:"userid" binding:"required"`
    PostId    uint      `gorm:"not null;uniqueIndex:idx_userid_star_postid,priority:10" json:"postid" binding:"required"`
}

type UserStarComment struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:,desc"`
    UserId    uint      `gorm:"not null;uniqueIndex:idx_userid_star_commentid,priority:8" json:"userid" binding:"required"`
    CommentId uint      `gorm:"not null;uniqueIndex:idx_userid_star_commentid,priority:10" json:"commentid" binding:"required"`
}

type UserAgreeComment struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:,desc"`
    UserId    uint      `gorm:"not null;uniqueIndex:idx_userid_agree_commentid,priority:8" json:"userid" binding:"required"`
    CommentId uint      `gorm:"not null;uniqueIndex:idx_userid_agree_commentid,priority:10" json:"commentid" binding:"required"`
}

// post form

type UserRegisterForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateForm struct {
	UserId      uint   `json:"userid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PasswordOld string `json:"passwordold" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
