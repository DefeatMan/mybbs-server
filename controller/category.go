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

func QueryCategorybyId(ctx *gin.Context) {
	categoryIdstr := ctx.Param("cid")
	categoryId, err := strconv.ParseInt(categoryIdstr, 10, 32)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "/:cid param error")
		return
	}
	category, err := dao.QueryCategorybyId(uint(categoryId))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category query failed")
		return
	}
	response.ResponseSuccess(ctx, gin.H{
		"id":     category.ID,
		"name":   category.Name,
		"follow": category.FollowNum,
	})
	return
}

func CreateCategory(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.CategoryPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	var ccform models.CategoryCreateForm
	if err := ctx.ShouldBindJSON(&ccform); err != nil {
		zap.L().Error("create with invalied params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}
	_, err = dao.CreateCategory(ccform.Name)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category create failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func RenameCategory(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.CategoryPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

	var rcform models.CategoryRenameForm
	if err := ctx.ShouldBindJSON(&rcform); err != nil {
		zap.L().Error("update with invalied params", zap.Error(err))
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(ctx, response.CodeInvalidParams)
			return
		}
		response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, err2)
		return
	}
	err = dao.RenameCategory(rcform.CatagoryId, rcform.Name)
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category rename failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func DeleteCategory(ctx *gin.Context) {
	egoId, exist := ctx.Get("egoId")
	if !exist {
		return
	}

	admin, err := dao.QueryAdminbyUserId(egoId.(uint))
	if err != nil || (admin.AdminPerm&models.CategoryPermFlag) == 0 {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "permission denied")
		return
	}

    categoryIdstr := ctx.Query("category")
    categoryId, err := strconv.ParseInt(categoryIdstr, 10, 32)
    if err != nil {
        response.ResponseErrorWithMsg(ctx, response.CodeInvalidParams, "?category= param error")
        return
    }

    err = dao.DeleteCategory(uint(categoryId))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category delete failed")
		return
	}
	response.ResponseSuccess(ctx, nil)
	return
}

func QueryCategoryList(ctx *gin.Context) {
	offsetStr := ctx.Query("offset")
	showStr := ctx.Query("show")
	offset_num, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset_num = 0
	}
	show_num, err := strconv.ParseInt(showStr, 10, 32)
	if err != nil {
		show_num = 20
	}

	categoryList, err := logic.CategoryListPage(int(offset_num), int(show_num))
	if err != nil {
		response.ResponseErrorWithMsg(ctx, response.CodeUnknownError, "category list query failed")
		return
	}
	response.ResponseSuccess(ctx, categoryList)
}
