package logic

import (
	"kome/mybbs-server/dao"
	"kome/mybbs-server/models"

	"github.com/gin-gonic/gin"
)

func categoryFill(categorys []models.Category) (categoryData []*gin.H, count int) {
	count = len(categorys)
	categoryData = make([]*gin.H, 0, count)
	for i := 0; i < count; i++ {
		categoryData = append(categoryData, &gin.H{
			"id":     categorys[i].ID,
			"name":   categorys[i].Name,
			"follow": categorys[i].FollowNum,
		})
	}
	return
}

func CategoryListPage(start_offset int, show_num int) (result *gin.H, err error) {
	total, err := dao.CountCategory()
	if err != nil {
		return
	}
	categorys, err := dao.QueryCategoryPage(start_offset, show_num)
	if err != nil {
		return
	}
	categoryData, count := categoryFill(categorys)
	result = &gin.H{
		"tot":        total,
		"cur":        count,
		"resultlist": categoryData,
	}
	return
}
