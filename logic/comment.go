package logic

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/models"
	"time"

	"github.com/gin-gonic/gin"
)

func commentOnlyFill(comments []models.Comment) (commentData []*gin.H, count int) {
	count = len(comments)
	commentData = make([]*gin.H, 0, count)
	for i := 0; i < count; i++ {
		commentData = append(commentData, &gin.H{
			"id":         comments[i].ID,
			"content":    comments[i].Content,
			"linkid":     comments[i].LinkId,
			"createtime": comments[i].CreatedAt,
            "star":comments[i].StarNum,
            "agree":comments[i].AgreeNum,
			"user": &gin.H{
				"id": comments[i].UserId,
			},
		})
	}
	return
}

func commentFill(comments []models.Comment, users []models.User) (commentData []*gin.H, count int) {
	count = len(comments)
	commentData = make([]*gin.H, 0, count)
	for i := 0; i < count; i++ {
		commentData = append(commentData, &gin.H{
			"id":         comments[i].ID,
			"content":    comments[i].Content,
			"linkid":     comments[i].LinkId,
			"createtime": comments[i].CreatedAt,
            "star" : comments[i].StarNum,
            "agree": comments[i].AgreeNum,
			"user": &gin.H{
				"id":    users[i].ID,
				"name":  users[i].Name,
				"agree": users[i].AgreeNum,
				"star":  users[i].StarNum,
			},
		})
	}
	return
}

func CommentRelated(commentId uint, show_num int) (result *gin.H, err error) {
	comments, users, err := dao.QueryCommentRelated(commentId, show_num)
	if err != nil {
		return
	}
	commentData, count := commentFill(comments, users)
	result = &gin.H{
		"cur":        count,
		"resultlist": commentData,
	}
	return
}

func CommentListbyPostIdPage(postId uint, start_offset int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountCommentbyPostId(postId)
	if err != nil {
		return
	}
	comments, users, end_time, err := dao.QueryCommentbyPostIdPage(postId, start_offset, show_num)
	if err != nil {
		return
	}
	commentData, count := commentFill(comments, users)
	result = &gin.H{
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": commentData,
	}
	return
}

func CommentListbyPostIdCreateTime(postId uint, start_time time.Time, show_num int) (result *gin.H, err error) {
	total, err := dao.CountCommentbyPostId(postId)
	if err != nil {
		return
	}
	comments, users, end_time, err := dao.QueryCommentbyPostIdCreateTime(postId, start_time, show_num)
	if err != nil {
		return
	}
	commentData, count := commentFill(comments, users)
	result = &gin.H{
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": commentData,
	}
	return
}

func CommentListbyPostIdAgreePage(postId uint, start_offset int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountCommentbyPostId(postId)
	if err != nil {
		return
	}
	comments, users, end_agreenum, end_time, err := dao.QueryCommentbyPostIdAgreePage(postId, start_offset, show_num)
	if err != nil {
		return
	}
	commentData, count := commentFill(comments, users)
	result = &gin.H{
		"endagree":   end_agreenum,
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": commentData,
	}
	return
}

func CommentListbyPostIdAgreeCreateTime(postId uint, start_agreenum uint, start_time time.Time, show_num int) (result *gin.H, err error) {
	total, err := dao.CountCommentbyPostId(postId)
	if err != nil {
		return
	}
	comments, users, end_agreenum, end_time, err := dao.QueryCommentbyPostIdAgreeCreateTime(postId, start_agreenum, start_time, show_num)
	if err != nil {
		return
	}
	commentData, count := commentFill(comments, users)
	result = &gin.H{
		"endagree":   end_agreenum,
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": commentData,
	}
	return
}
