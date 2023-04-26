package middle

import (
	"kome/mybbs-server/response"
	"kome/mybbs-server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func EgoMiddle() func(*gin.Context) {
	return func(ctx *gin.Context) {
		// token : header.Authorization [Bearer ***]
		userHeader := ctx.Request.Header.Get("Authorization")
		if userHeader == "" {
			ctx.Next()
			return
		}
		parts := strings.SplitN(userHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.Next()
			return
		}
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			ctx.Next()
			return
		}
		ctx.Set("egoId", mc.UserId)
		ctx.Next()
	}
}

func NeedEgoMiddle() func(*gin.Context) {
	return func(ctx *gin.Context) {
		// token : header.Authorization [Bearer ***]
		userHeader := ctx.Request.Header.Get("Authorization")
		if userHeader == "" {
			response.ResponseErrorWithMsg(ctx, response.CodeInvalidAuthFormat, "header authorization format err")
			ctx.Abort()
			return
		}
		parts := strings.SplitN(userHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.ResponseErrorWithMsg(ctx, response.CodeInvalidAuthFormat, "header authorization format err")
			ctx.Abort()
			return
		}
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			response.ResponseError(ctx, response.CodeInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("egoId", mc.UserId)
		ctx.Next()
	}
}
