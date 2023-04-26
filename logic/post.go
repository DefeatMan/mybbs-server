package logic

import (
	"github.com/gin-gonic/gin"
	"kome/mybbs-server/dao"
	"kome/mybbs-server/models"
	"time"
)

func postFill(posts []models.Post) (postData []*gin.H, count int) {
	count = len(posts)
	postData = make([]*gin.H, 0, count)
	for i := 0; i < count; i++ {
		postData = append(postData, &gin.H{
			"id":         posts[i].ID,
			"title":      posts[i].Title,
			"createtime": posts[i].CreatedAt,
			"commentid":  posts[i].CommentId,
			"category":   posts[i].CategoryId,
			"star":       posts[i].StarNum,
			"lock":       posts[i].LockFlag,
			"user":       posts[i].UserId,
		})
	}
	return
}

func PostListbyCategoryIdPage(categoryId uint, start_offset int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountPostbyCategoryId(categoryId)
	if err != nil {
		return
	}
	posts, end_time, err := dao.QueryPostbyCategoryIdPage(categoryId, start_offset, show_num)
	if err != nil {
		return
	}
	postData, count := postFill(posts)
	result = &gin.H{
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": postData,
	}
	return
}

func PostListbyCategoryIdCreateTime(categoryId uint, start_time time.Time, show_num int) (result *gin.H, err error) {
	total, err := dao.CountPostbyCategoryId(categoryId)
	if err != nil {
		return
	}
	posts, end_time, err := dao.QueryPostbyCategoryIdCreateTime(categoryId, start_time, show_num)
	if err != nil {
		return
	}
	postData, count := postFill(posts)
	result = &gin.H{
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": postData,
	}
	return
}

func PostListbyCategoryIdStarPage(categoryId uint, start_offset int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountPostbyCategoryId(categoryId)
	if err != nil {
		return
	}
	posts, end_starnum, end_time, err := dao.QueryPostbyCategoryIdStarPage(categoryId, start_offset, show_num)
	if err != nil {
		return
	}
	postData, count := postFill(posts)
	result = &gin.H{
		"endstar":    end_starnum,
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": postData,
	}
	return
}

func PostListbyCategoryIdStarCreateTime(categoryId uint, start_starnum uint, start_time time.Time, show_num int) (result *gin.H, err error) {
	total, err := dao.CountPostbyCategoryId(categoryId)
	if err != nil {
		return
	}
	posts, end_starnum, end_time, err := dao.QueryPostbyCategoryIdStarCreateTime(categoryId, start_starnum, start_time, show_num)
	if err != nil {
		return
	}
	postData, count := postFill(posts)
	result = &gin.H{
		"end_star":   end_starnum,
		"endtime":    end_time,
		"tot":        total,
		"cur":        count,
		"resultlist": postData,
	}
	return
}
