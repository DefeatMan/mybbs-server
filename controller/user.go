package controller

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/models"
	"kome/mybbs-server/response"
	"kome/mybbs-server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// user register
func UserRegister(ctx *gin.Context) {
	var urform models.UserRegisterForm
	if err := ctx.ShouldBindJSON(&urform); err != nil {
		zap.L().Error("sign up with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			// return invalidparams
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		// return error detail
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}
	_, err := dao.CreateUser(urform.Name, urform.Email, urform.Password)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "user register failed: "+err.Error())
		return
	}
	response.ResponseSuccess(ctx, nil)
}

// user login
func UserLogin(ctx *gin.Context) {
	var ulform models.UserLoginForm
	if err := ctx.ShouldBindJSON(&ulform); err != nil {
		zap.L().Error("sign in with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			// return invalidparams
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		// return error detail
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}
	user, err := dao.UserLogin(ulform.Email, ulform.Password)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "user login failed")
		return
	}

	token, err := utils.NewToken(user.ID, user.Email)

	response.ResponseSuccess(ctx, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"agree": user.AgreeNum,
		"star":  user.StarNum,
		"token": token,
	})
}

func QueryUserbyId(ctx *gin.Context) {
	userIdstr := ctx.Param("id")
	userId, err := strconv.ParseInt(userIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "/:id param error")
		return
	}

	user, err := dao.QueryUserbyId(uint(userId))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "user query failed")
		return
	}

	admin, err := dao.QueryAdminbyUserId(uint(userId))
	var adminPerm uint = 0
	if err == nil {
		adminPerm = admin.AdminPerm
	}

	response.ResponseSuccess(ctx, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"agree": user.AgreeNum,
		"star":  user.StarNum,
		"admin": adminPerm,
	})
	return
}

func UpdateUser(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}
	var uuform models.UserUpdateForm
	if err := ctx.ShouldBindJSON(&uuform); err != nil {
		zap.L().Error("update user with invalid params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}

	userId := uuform.UserId
	if uint(userId) != egoId.(uint) {
		egoadmin, err := dao.QueryAdminbyUserId(egoId.(uint))
		if err != nil || (egoadmin.AdminPerm&models.AdminPermFlag) == 0 {
			response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
			return
		}
	}

	_, err := dao.UpdateUser(uuform.UserId, uuform.Name, uuform.Email, uuform.Password, uuform.PasswordOld)

	if err != nil {
		zap.L().Error("user update failed", zap.Error(err))
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "user update failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func DeleteUser(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	userIdstr := ctx.Query("id")
	userId, err := strconv.ParseInt(userIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "?id= param error")
		return
	}

	if uint(userId) != egoId.(uint) {
		egoadmin, err := dao.QueryAdminbyUserId(egoId.(uint))
		if err != nil || (egoadmin.AdminPerm&models.AdminPermFlag) == 0 {
			response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
			return
		}
	}

	err = dao.DeleteUser(uint(userId))

	if err != nil {
		zap.L().Error("user delete failed", zap.Error(err))
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "user delete failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}
