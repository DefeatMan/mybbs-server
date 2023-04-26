package router

import (
	"kome/mybbs-server/controller"
	"kome/mybbs-server/middle"

	"github.com/gin-gonic/gin"
)

func userRoute(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/cadmin", middle.NeedEgoMiddle(), controller.CreateAdmin)
	r.PUT("/aadmin", middle.NeedEgoMiddle(), controller.AppendAdmin)
	r.PUT("/sadmin", middle.NeedEgoMiddle(), controller.SubsetAdmin)
	r.DELETE("/dadmin", middle.NeedEgoMiddle(), controller.DeleteAdmin)

	r.POST("/register", controller.UserRegister)
	r.POST("/login", controller.UserLogin)
	r.PUT("/uuser", middle.NeedEgoMiddle(), controller.UpdateUser)
	r.DELETE("/duser", middle.NeedEgoMiddle(), controller.DeleteUser)

	r.POST("/ccategory", middle.NeedEgoMiddle(), controller.CreateCategory)
	r.POST("/ucategory", middle.NeedEgoMiddle(), controller.RenameCategory)
	r.DELETE("/dcategory", middle.NeedEgoMiddle(), controller.DeleteCategory)

	r.POST("/cpost", middle.NeedEgoMiddle(), controller.CreatePost)
	r.DELETE("/dpost", middle.NeedEgoMiddle(), controller.DeletePost)
	r.DELETE("/lpost", middle.NeedEgoMiddle(), controller.LockPost)

	r.POST("/ccomment", middle.NeedEgoMiddle(), controller.CreateComment)
	r.PUT("/acomment", middle.NeedEgoMiddle(), controller.AppendComment)
	r.DELETE("/dcomment", middle.NeedEgoMiddle(), controller.DeleteComment)

	return r
}

func userStarRoute(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/post/:pid", middle.NeedEgoMiddle(), controller.ClickStarPost)
	r.POST("/comment/:cid", middle.NeedEgoMiddle(), controller.ClickStarComment)
	r.GET("/post", middle.NeedEgoMiddle(), controller.QueryStarPost)
	r.GET("/comment", middle.NeedEgoMiddle(), controller.QueryStarComment)
    
	return r
}

func userAgreeRoute(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/comment/:cid", middle.NeedEgoMiddle(), controller.ClickAgreeComment)

	return r
}

func queryRoute(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/category", controller.QueryCategoryList)
	r.GET("/category/:cid", controller.QueryPostList)
	r.GET("/category/:cid/star", controller.QueryPostListStar)
	r.GET("/post/:pid", controller.QueryCommentList)
	r.GET("/post/:pid/agree", controller.QueryCommentListAgree)

	return r
}

func singleRoute(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/user/:id", controller.QueryUserbyId)
    r.GET("/category/:cid", controller.QueryCategorybyId)
	r.GET("/post/:pid", controller.QueryPostbyId)
    r.GET("/comment/:cid", controller.QueryCommentbyId)
	r.GET("/comment/:cid/related", controller.QueryCommentRelated)

	return r
}

func InitRoute(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	s := api.Group("/s")
	s = singleRoute(s)
    star := api.Group("/star")
	star = userStarRoute(star)
	agree := api.Group("/agree")
	agree = userAgreeRoute(agree)
	api = userRoute(api)
	api = queryRoute(api)

	return r
}
