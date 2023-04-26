package controller

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/models"
	"kome/mybbs-server/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func CreateAdmin(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.AdminPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	var acform *models.Admin
	if err := ctx.ShouldBindJSON(acform); err != nil {
		zap.L().Error("create with invalied params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}
	err = dao.CreateAdmin(acform.UserId, acform.AdminPerm&admin.AdminPerm)

	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "admin create failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func AppendAdmin(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.AdminPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	var acform *models.Admin
	if err := ctx.ShouldBindJSON(acform); err != nil {
		zap.L().Error("create with invalied params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	err = dao.AppendAdminPerm(acform.UserId, acform.AdminPerm&admin.AdminPerm)

	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "append admin perm failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func SubsetAdmin(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.AdminPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	var acform *models.Admin
	if err := ctx.ShouldBindJSON(acform); err != nil {
		zap.L().Error("create with invalied params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	err = dao.AppendAdminPerm(acform.UserId, acform.AdminPerm|(^admin.AdminPerm))

	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "subset admin perm failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func DeleteAdmin(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	userIdstr := ctx.Query("id")
	userId, err := strconv.ParseInt(userIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "?id= param error")
		return
	}

	if userId != egoId && admin.AdminPerm != models.RootPermFlag {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}
	err = dao.DeleteAdmin(uint(userId))

	if err != nil {
		zap.L().Error("user delete failed", zap.Error(err))
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "admin delete failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}
