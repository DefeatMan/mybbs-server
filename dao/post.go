package dao

import (
	"gorm.io/gorm"
	"kome/mybbs-server/models"
	"time"
)

func CountPostbyCategoryId(categoryId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.Post{}).Where("category_id = ?", categoryId).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func QueryPostbyId(postId uint) (post *models.Post, err error) {
	db := GetDatabase()
	post = new(models.Post)
	result := db.Where("id = ?", postId).Limit(1).Find(post)
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

func QueryPostbyCategoryIdCreateTime(categoryId uint, start_time time.Time, show_num int) (post_list []models.Post, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("category_id = ? AND created_at < ?", categoryId, start_time).Order("created_at desc").Limit(show_num).Find(&post_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_time = post_list[len(post_list)-1].CreatedAt
	return
}

func QueryPostbyCategoryIdPage(categoryId uint, offset_num int, show_num int) (post_list []models.Post, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("category_id = ?", categoryId).Order("created_at desc").Offset(offset_num).Limit(show_num).Find(&post_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_time = post_list[len(post_list)-1].CreatedAt
	return
}

func QueryPostbyCategoryIdStarCreateTime(categoryId uint, start_starnum uint, start_time time.Time, show_num int) (post_list []models.Post, end_starnum uint, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("category_id = ? AND star_num = ? AND created_at < ?", categoryId, start_starnum, start_time).Order("created_at desc").Limit(show_num).Find(&post_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	show_num -= len(post_list)
	if show_num > 0 {
		var post_tail []models.Post
		result = db.Where("category_id = ? AND star_num < ?", categoryId, start_starnum).Order("star_num desc, created_at desc").Limit(show_num).Find(&post_tail)
		if result.Error != nil {
			err = EQueryFailed
			return
		}
		for i := range post_tail {
			post_list = append(post_list, post_tail[i])
		}
	}
	if len(post_list) == 0 {
		err = ENotExist
		return
	}

	end_starnum = post_list[len(post_list)-1].StarNum
	end_time = post_list[len(post_list)-1].CreatedAt
	return
}

func QueryPostbyCategoryIdStarPage(categoryId uint, offset_num int, show_num int) (post_list []models.Post, end_starnum uint, end_time time.Time, err error) {
	db := GetDatabase()
	result := db.Where("category_id = ?", categoryId).Order("star_num desc, created_at desc").Offset(offset_num).Limit(show_num).Find(&post_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	end_starnum = post_list[len(post_list)-1].StarNum
	end_time = post_list[len(post_list)-1].CreatedAt
	return
}

func CreatePost(categoryId uint, title string, content string, userId uint) (post *models.Post, comment *models.Comment, err error) {
	db := GetDatabase()
	err = db.Transaction(func(tx *gorm.DB) error {
		comment = &models.Comment{
			PostId:   models.BaseTextFlag,
			LinkId:   models.BaseTextFlag,
			UserId:   userId,
			Content:  content,
			AgreeNum: 0,
		}
		if tx.Create(comment).Error != nil {
			return ECreateFailed
		}

		post = &models.Post{
			CategoryId: categoryId,
			Title:      title,
			UserId:     userId,
			CommentId:  comment.ID,
			StarNum:    0,
			LockFlag:   0,
		}
		if tx.Create(post).Error != nil {
			return ECreateFailed
		}
		if tx.Model(&models.Category{}).Where("id = ?", categoryId).Update("follow_num", gorm.Expr("follow_num + ?", 1)).Error != nil {
			return EUpdateFailed
		}
		return nil
	})
	return
}

func LockPost(postId uint, userId uint, adminPerm uint) error {
	db := GetDatabase()
	if db.Model(&models.Post{PostModel: models.PostModel{ID: postId}}).Where("user_id = ? OR 0 < ?", userId, adminPerm).Update("lock_flag", 1).Error != nil {
		return EUpdateFailed
	}
	return nil
}

func DeletePost(postId uint, userId uint, adminPerm uint) error {
	db := GetDatabase()
	if db.Where("user_id = ? OR 0 < ?", userId, adminPerm).Delete(&models.Post{}, postId).Error != nil {
		return EDeleteFailed
	}
	return nil
}
