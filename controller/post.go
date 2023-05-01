package controller

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/logic"
	"kome/mybbs-server/models"
	"kome/mybbs-server/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func QueryPostbyId(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	postIdstr := ctx.Param("pid")
	postId, err := strconv.ParseInt(postIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "/:pid param error")
		return
	}
	post, err := dao.QueryPostbyId(uint(postId))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category query failed")
		return
	}
	var star int64 = 0
	if exist {
		star, err = dao.StarPostCheck(egoId.(uint), uint(postId))
	}
	response.ResponseSuccess(ctx, gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"createtime": post.CreatedAt,
		"commentid":  post.CommentId,
		"category":   post.CategoryId,
		"star":       post.StarNum,
		"lock":       post.LockFlag,
		"user":       post.UserId,
		"mystar":     star,
	})
	return
}

func CreatePost(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	var cpform models.PostCreateForm
	if err := ctx.ShouldBindJSON(&cpform); err != nil {
		zap.L().Error("create post with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	if _, _, err := dao.CreatePost(cpform.CategoryId, cpform.Title, cpform.Content, egoId.(uint)); err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "post create failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func LockPost(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	var adminPerm uint = 0
	if err == nil {
		adminPerm = admin.AdminPerm
	}

	postIdstr := ctx.Query("post")
	postId, err := strconv.ParseInt(postIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "?post= param error")
		return
	}
	err = dao.LockPost(uint(postId), egoId.(uint), adminPerm & models.PostPermFlag)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "post lock failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func DeletePost(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	var adminPerm uint = 0
	if err == nil {
		adminPerm = admin.AdminPerm
	}

	postIdstr := ctx.Query("post")
	postId, err := strconv.ParseInt(postIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "?post= param error")
		return
	}
	err = dao.DeletePost(uint(postId), egoId.(uint), adminPerm & models.PostPermFlag)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "post delete failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func QueryPostList(ctx *gin.Context) {
	categoryStr := ctx.Param("cid")
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	categoryId, err := strconv.ParseInt(categoryStr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "category id parse error")
		return
	}
	offset_num, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset_num = 0
	}
	show_num, err := strconv.ParseInt(showStr, 10, 32)
	if err != nil {
		show_num = 20
	}

	postList, err := logic.PostListbyCategoryIdPage(uint(categoryId), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "post list query failed")
		return
	}
	response.ResponseSuccess(ctx, postList)
	return
}

func QueryPostListStar(ctx *gin.Context) {
	categoryStr := ctx.Param("cid")
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	categoryId, err := strconv.ParseInt(categoryStr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "category id parse error")
		return
	}
	offset_num, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset_num = 0
	}
	show_num, err := strconv.ParseInt(showStr, 10, 32)
	if err != nil {
		show_num = 20
	}

	postList, err := logic.PostListbyCategoryIdStarPage(uint(categoryId), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "post list query failed")
		return
	}
	response.ResponseSuccess(ctx, postList)
	return
}
