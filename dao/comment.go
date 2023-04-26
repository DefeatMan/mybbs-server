package dao

import (
	"kome/mybbs-server/models"
	"strings"
	"time"
)

func CountCommentbyPostId(postId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.Comment{}).Where("post_id = ?", postId).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func QueryCommentbyId(commentId uint) (res *models.Comment, err error) {
	db := GetDatabase()
	res = new(models.Comment)
	result := db.Where("id = ?", commentId).Limit(1).Find(res)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}
	return
}

func QueryCommentbyPostIdCreateTime(postId uint, start_time time.Time, show_num int) (text_list []models.Comment, user_list []models.User, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("post_id = ? AND created_at > ?", postId, start_time).Order("created_at asc").Limit(show_num).Find(&text_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_time = text_list[len(text_list)-1].CreatedAt

	for _, text_one := range text_list {
		user_one, err2 := QueryUserbyId(text_one.UserId)
		if err2 != nil {
			err = err2
			return
		}
		user_list = append(user_list, *user_one)
	}
	return
}

func QueryCommentbyPostIdPage(postId uint, offset_num int, show_num int) (text_list []models.Comment, user_list []models.User, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("post_id = ?", postId).Order("created_at asc").Offset(offset_num).Limit(show_num).Find(&text_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_time = text_list[len(text_list)-1].CreatedAt

	for _, text_one := range text_list {
		user_one, err2 := QueryUserbyId(text_one.UserId)
		if err2 != nil {
			err = err2
			return
		}
		user_list = append(user_list, *user_one)
	}
	return
}

func QueryCommentbyPostIdAgreeCreateTime(postId uint, start_agreenum uint, start_time time.Time, show_num int) (text_list []models.Comment, user_list []models.User, end_agreenum uint, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("post_id = ? AND agree_num = ? AND created_at > ?", postId, start_agreenum, start_time).Order("created_at asc").Limit(show_num).Find(&text_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	show_num -= len(text_list)
	if show_num > 0 {
		var text_tail []models.Comment
		result = db.Where("post_id = ? AND agree_num < ?", postId, start_agreenum).Order("agree_num desc, created_at asc").Limit(show_num).Find(&text_tail)
		if result.Error != nil {
			err = EQueryFailed
			return
		}
		for i := range text_tail {
			text_list = append(text_list, text_tail[i])
		}
	}
	if len(text_list) == 0 {
		err = ENotExist
		return
	}

	end_agreenum = text_list[len(text_list)-1].AgreeNum
	end_time = text_list[len(text_list)-1].CreatedAt

	for i := range text_list {
		user_one, err2 := QueryUserbyId(text_list[i].UserId)
		if err2 != nil {
			err = err2
			return
		}
		user_list = append(user_list, *user_one)
	}
	return
}

func QueryCommentbyPostIdAgreePage(postId uint, offset_num int, show_num int) (text_list []models.Comment, user_list []models.User, end_agreenum uint, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("post_id = ?", postId).Order("agree_num desc, created_at asc").Offset(offset_num).Limit(show_num).Find(&text_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_agreenum = text_list[len(text_list)-1].AgreeNum
	end_time = text_list[len(text_list)-1].CreatedAt

	for i := range text_list {
		user_one, err2 := QueryUserbyId(text_list[i].UserId)
		if err2 != nil {
			err = err2
			return
		}
		user_list = append(user_list, *user_one)
	}
	return
}

func CreateComment(postId uint, linkId uint, userId uint, content string) (comment *models.Comment, err error) {
	comment = &models.Comment{
		PostId:   postId,
		LinkId:   linkId,
		UserId:   userId,
		Content:  content,
		AgreeNum: 0,
		StarNum:  0,
	}

	db := GetDatabase()
	if db.Create(comment).Error != nil {
		err = ECreateFailed
		return
	}
	return
}

func AppendComment(userId uint, commentId uint, content string) (comment *models.Comment, err error) {
	db := GetDatabase()
	comment = new(models.Comment)
	result := db.Where("id = ? AND user_id = ?", commentId, userId).Limit(1).Find(comment)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}
	var builder strings.Builder
	builder.WriteString(comment.Content)
	builder.Write([]byte(`<br> *---NOTICE: append content ---* <br>`))
	builder.WriteString(content)
	comment.Content = builder.String()

	if db.Model(comment).Update("content", comment.Content).Error != nil {
		err = EUpdateFailed
		return
	}
	return
}

func DeleteComment(commentId uint, userId uint, adminPerm uint) error {
	db := GetDatabase()
	if db.Model(&models.Comment{}).Where("id = ? AND (user_id = ? OR 0 < ?)", commentId, userId, adminPerm).Update("content", "comment had been deleted").Error != nil {
		return EDeleteFailed
	}
	return nil
}

func QueryCommentRelated(commentId uint, show_num int) (text_list []models.Comment, user_list []models.User, err error) {
	CatchOne := func(commentId uint) (parentId uint, err error) {
		text_one, err := QueryCommentbyId(commentId)
		if err != nil {
			return
		}
		user_one, err := QueryUserbyId(text_one.UserId)
		if err != nil {
			return
		}
		text_list = append(text_list, *text_one)
		user_list = append(user_list, *user_one)
		return
	}
	now_commentId := commentId
	for i := 0; i < show_num; i++ {
		now_commentId, err = CatchOne(now_commentId)
		if err != nil {
			return
		}
		if now_commentId == 0 {
			return
		}
	}
	return
}
