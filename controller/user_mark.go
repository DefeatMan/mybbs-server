package controller

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryStarPost(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	offset_num, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset_num = 10
	}
	show_num, err := strconv.ParseInt(showStr, 10, 32)
	if err != nil {
		show_num = 20
	}

	postList, err := dao.QueryStarPostList(egoId.(uint), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "star post list query failed")
		return
	}
	response.ResponseSuccess(ctx, postList)
	return
}

func QueryStarComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	offset_num, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset_num = 10
	}
	show_num, err := strconv.ParseInt(showStr, 10, 32)
	if err != nil {
		show_num = 20
	}

	commentList, err := dao.QueryStarCommentList(egoId.(uint), int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "star comment list query failed")
		return
	}
	response.ResponseSuccess(ctx, commentList)
	return
}

func ClickStarPost(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	idStr := ctx.Param("pid")
	id, err := strconv.ParseInt(idStr, 10, 32)
	count, now_state, err := dao.ClickStarPost(egoId.(uint), uint(id))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "click star post failed")
		return
	}
	response.ResponseSuccess(ctx, gin.H{
		"count": count,
		"now":   now_state,
	})
	return
}

func ClickStarComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	idStr := ctx.Param("cid")
	id, err := strconv.ParseInt(idStr, 10, 32)
	count, now_state, err := dao.ClickStarComment(egoId.(uint), uint(id))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "click star comment failed")
		return
	}
	response.ResponseSuccess(ctx, gin.H{
		"count": count,
		"now":   now_state,
	})
	return
}

func ClickAgreeComment(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	idStr := ctx.Param("cid")
	id, err := strconv.ParseInt(idStr, 10, 32)
	count, now_state, err := dao.ClickAgreeComment(egoId.(uint), uint(id))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "click agree comment failed")
		return
	}
	response.ResponseSuccess(ctx, gin.H{
		"count": count,
		"now":   now_state,
	})
	return
}
