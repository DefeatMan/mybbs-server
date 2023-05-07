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

func QueryCommentbyId(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	commentIdstr := ctx.Param("cid")
	commentId, err := strconv.ParseInt(commentIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "comment param error")
		return
	}
	comment, err := dao.QueryCommentbyId(uint(commentId))

	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "comment query failed")
		return
	}

	user, err := dao.QueryUserbyId(comment.UserId)
	if err == dao.EQueryFailed {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "user query failed")
		return
	}

	var star int64 = 0
	var agree int64 = 0
	if exist {
		star, err = dao.StarCommentCheck(egoId.(uint), uint(commentId))
		agree, err = dao.AgreeCommentCheck(egoId.(uint), uint(commentId))
	}
	response.ResponseSuccess(ctx, gin.H{
		"id":         comment.ID,
		"content":    comment.Content,
		"linkid":       comment.LinkId,
		"createtime": comment.CreatedAt,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"agree": user.AgreeNum,
			"star":  user.StarNum,
		},
		"mystar":  star,
		"myagree": agree,
	})
	return
}

func CreateComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	var ccform models.Comment
	if err := ctx.ShouldBindJSON(&ccform); err != nil {
		zap.L().Error("create comment with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	if _, err := dao.CreateComment(ccform.PostId, ccform.LinkId, egoId.(uint), ccform.Content); err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment create failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func AppendComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	var acform models.CommentAppendForm
	if err := ctx.ShouldBindJSON(&acform); err != nil {
		zap.L().Error("append comment with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	if _, err := dao.AppendComment(egoId.(uint), acform.CommentId, acform.Content); err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment append failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func DeleteComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	var adminPerm uint = 0
	if err == nil {
		adminPerm = admin.AdminPerm
	}

	commentIdstr := ctx.Query("comment")
	commentId, err := strconv.ParseInt(commentIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "comment param error")
		return
	}
	err = dao.DeleteComment(uint(commentId), egoId.(uint), adminPerm & models.CommentPermFlag)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment delete failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func QueryCommentRelated(ctx *gin.Context) {
	commentStr := ctx.Param("cid")
	show_numStr := ctx.Query("show")
	commentId, err := strconv.ParseInt(commentStr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "comment id parse error")
		return
	}
	show_num, err := strconv.ParseInt(show_numStr, 10, 32)
	if err != nil {
		show_num = 5
	}

	commentList, err := logic.CommentRelated(uint(commentId), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment related query failed")
		return
	}
	response.ResponseSuccess(ctx, commentList)
}

func QueryCommentList(ctx *gin.Context) {
	postStr := ctx.Param("pid")
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	postId, err := strconv.ParseInt(postStr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "post id parse error")
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

	commentList, err := logic.CommentListbyPostIdPage(uint(postId), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment list query failed")
		return
	}
	response.ResponseSuccess(ctx, commentList)
	return
}

func QueryCommentListAgree(ctx *gin.Context) {
	postStr := ctx.Param("pid")
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	postId, err := strconv.ParseInt(postStr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "post id parse error")
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

	commentList, err := logic.CommentListbyPostIdAgreePage(uint(postId), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "comment list query failed")
		return
	}
	response.ResponseSuccess(ctx, commentList)
	return
}
