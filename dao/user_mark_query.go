package dao

import "kome/mybbs-server/models"

func QueryStarPostList(userId uint, offset_num int, show_num int) (posts []models.Post, err error) {
	db := GetDatabase()
	result := db.Raw("select posts.* from posts right join (select post_id, created_at from user_star_posts where user_id = ? order by created_at DESC limit ? offset ?) as tmp on posts.id = tmp.post_id order by tmp.created_at DESC", userId, show_num, offset_num).Scan(&posts)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func QueryStarCommentList(userId uint, offset_num int, show_num int) (comments []models.Comment, err error) {
	db := GetDatabase()
	result := db.Raw("select comments.* from comments right join (select comment_id, created_at from user_star_comments where user_id = ? order by created_at DESC limit ? offset ?) as tmp on comments.id = tmp.comment_id order by tmp.created_at DESC", userId, show_num, offset_num).Scan(&comments)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func CountStarPostbyUserId(userId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.UserStarPost{}).Where("user_id = ?", userId).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func CountStarCommentbyUserId(userId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.UserStarComment{}).Where("user_id = ?", userId).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}
