package dao

import (
	"kome/mybbs-server/models"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

func StarPostCheck(userId uint, postId uint) (count int64, err error) {
	db := GetDatabase()
	_ = db.Model(&models.UserStarPost{}).Where("user_id = ? AND post_id = ?", userId, postId).Limit(1).Count(&count)
	return
}

func StarCommentCheck(userId uint, commentId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.UserStarComment{}).Where("user_id = ? AND comment_id = ?", userId, commentId).Limit(1).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func AgreeCommentCheck(userId uint, commentId uint) (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.UserAgreeComment{}).Where("user_id = ? AND comment_id = ?", userId, commentId).Limit(1).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func ClickStarPost(userId uint, postId uint) (star_count uint, now_state uint, err error) {
	db := GetDatabase()
	count, err := StarPostCheck(userId, postId)
	if err != nil {
		return
	}

	if count == 0 {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 0
			user_star_post := &models.UserStarPost{
				UserId: userId,
				PostId: postId,
			}
			if tx.Create(user_star_post).Error != nil {
				return ECreateFailed
			}
			//var tmp_post []models.Post
			//result := tx.Model(&tmp_post).Clauses(clause.Returning{Columns: []clause.Column{{Name: "star_num"}, {Name: "user_id"}}}).Where("id = ?", postId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			result := tx.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_post models.Post
			result = tx.First(&tmp_post, postId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//star_count = tmp_post[0].StarNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_post[0].UserId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			star_count = tmp_post.StarNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_post.UserId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 1
			return nil
		})
	} else {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 1
			if tx.Unscoped().Where("user_id = ? AND post_id = ?", userId, postId).Delete(&models.UserStarPost{}).Error != nil {
				return EDeleteFailed
			}
			//var tmp_post []models.Post
			//result := tx.Model(&tmp_post).Clauses(clause.Returning{Columns: []clause.Column{{Name: "star_num"}, {Name: "user_id"}}}).Where("id = ?", postId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			result := tx.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_post models.Post
			result = tx.First(&tmp_post, postId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//star_count = tmp_post[0].StarNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_post[0].UserId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			star_count = tmp_post.StarNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_post.UserId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 0
			return nil
		})
	}
	return
}

func ClickStarComment(userId uint, commentId uint) (star_count uint, now_state uint, err error) {
	db := GetDatabase()
	count, err := StarCommentCheck(userId, commentId)
	if err != nil {
		return
	}

	if count == 0 {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 0
			user_star_comment := &models.UserStarComment{
				UserId:    userId,
				CommentId: commentId,
			}
			if tx.Create(user_star_comment).Error != nil {
				return ECreateFailed
			}
			//var tmp_comment []models.Comment
			//result := tx.Model(&tmp_comment).Clauses(clause.Returning{Columns: []clause.Column{{Name: "star_num"}, {Name: "user_id"}}}).Where("id = ?", commentId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			result := tx.Model(&models.Comment{}).Where("id = ?", commentId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_comment models.Comment
			result = tx.First(&tmp_comment, commentId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//star_count = tmp_comment[0].StarNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_comment[0].UserId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			star_count = tmp_comment.StarNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_comment.UserId).UpdateColumn("star_num", gorm.Expr("star_num + ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 1
			return nil
		})
	} else {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 1
			if tx.Unscoped().Where("user_id = ? comment_id = ?", userId, commentId).Delete(&models.UserStarComment{}).Error != nil {
				return EDeleteFailed
			}
			//var tmp_comment []models.Comment
			//result := tx.Model(&tmp_comment).Clauses(clause.Returning{Columns: []clause.Column{{Name: "star_num"}, {Name: "user_id"}}}).Where("id = ?", commentId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			result := tx.Model(&models.Comment{}).Where("id = ?", commentId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_comment models.Comment
			result = tx.First(&tmp_comment, commentId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//star_count = tmp_comment[0].StarNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_comment[0].UserId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			star_count = tmp_comment.StarNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_comment.UserId).UpdateColumn("star_num", gorm.Expr("star_num - ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 0
			return nil
		})
	}
	return
}

func ClickAgreeComment(userId uint, commentId uint) (agree_count uint, now_state uint, err error) {
	db := GetDatabase()
	count, err := AgreeCommentCheck(userId, commentId)
	if err != nil {
		return
	}

	if count == 0 {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 0
			user_agree_comment := &models.UserAgreeComment{
				UserId:    userId,
				CommentId: commentId,
			}
			if tx.Create(user_agree_comment).Error != nil {
				return ECreateFailed
			}
			//var tmp_comment []models.Comment
			//result := tx.Model(&tmp_comment).Clauses(clause.Returning{Columns: []clause.Column{{Name: "agree_num"}, {Name: "user_id"}}}).Where("id = ?", commentId).UpdateColumn("agree_num", gorm.Expr("agree_num + ?", 1))
			result := tx.Model(&models.Comment{}).Where("id = ?", commentId).UpdateColumn("agree_num", gorm.Expr("agree_num + ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_comment models.Comment
			result = tx.First(&tmp_comment, commentId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//agree_count = tmp_comment[0].AgreeNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_comment[0].UserId).UpdateColumn("agree_num", gorm.Expr("agree_num + ?", 1))
			agree_count = tmp_comment.AgreeNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_comment.UserId).UpdateColumn("agree_num", gorm.Expr("agree_num + ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 1
			return nil
		})
	} else {
		err = db.Transaction(func(tx *gorm.DB) error {
			now_state = 1
			if tx.Unscoped().Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&models.UserAgreeComment{}).Error != nil {
				return EDeleteFailed
			}
			//var tmp_comment []models.Comment
			//result := tx.Model(&tmp_comment).Clauses(clause.Returning{Columns: []clause.Column{{Name: "agree_num"}, {Name: "user_id"}}}).Where("id = ?", commentId).UpdateColumn("agree_num", gorm.Expr("agree_num - ?", 1))
			result := tx.Model(&models.Comment{}).Where("id = ?", commentId).UpdateColumn("agree_num", gorm.Expr("agree_num - ?", 1))
			if result.Error != nil || result.RowsAffected == 0 {
				return EUpdateFailed
			}
			var tmp_comment models.Comment
			result = tx.First(&tmp_comment, commentId)
			if result.Error != nil {
				return EUpdateFailed
			}
			//agree_count = tmp_comment[0].AgreeNum
			//result = tx.Model(&models.User{}).Where("id = ?", tmp_comment[0].UserId).UpdateColumn("agree_num", gorm.Expr("agree_num - ?", 1))
			agree_count = tmp_comment.AgreeNum
			result = tx.Model(&models.User{}).Where("id = ?", tmp_comment.UserId).UpdateColumn("agree_num", gorm.Expr("agree_num - ?", 1))
			if result.Error != nil {
				return EUpdateFailed
			}
			now_state = 0
			return nil
		})
	}
	return
}
