package logic

import (
	"kome/mybbs-server/dao"

	"github.com/gin-gonic/gin"
)

func UserStarPostList(userId uint, offset_num int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountStarPostbyUserId(userId)
	if err != nil {
		return
	}
	posts, err := dao.QueryStarPostList(userId, offset_num, show_num)
	if err != nil {
		return
	}
	postData, count := postFill(posts)
	result = &gin.H{
		"tot":        total,
		"cur":        count,
		"resultlist": postData,
	}
	return
}

func UserStarCommentList(userId uint, offset_num int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountStarCommentbyUserId(userId)
	if err != nil {
		return
	}
	comments, err := dao.QueryStarCommentList(userId, offset_num, show_num)
	if err != nil {
		return
	}
	commentData, count := commentOnlyFill(comments)
	result = &gin.H{
		"tot":        total,
		"cur":        count,
		"resultlist": commentData,
	}
	return
}
