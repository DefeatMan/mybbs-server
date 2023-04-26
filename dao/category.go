package dao

import "kome/mybbs-server/models"

func CountCategory() (count int64, err error) {
	db := GetDatabase()
	result := db.Model(&models.Category{}).Count(&count)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	return
}

func QueryCategorybyId(categoryId uint) (category *models.Category, err error) {
	db := GetDatabase()
	category = new(models.Category)
	result := db.Where("id = ?", categoryId).Limit(1).Find(category)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}
	return
}

func QueryCategoryPage(offset_num int, show_num int) (category_list []models.Category, err error) {
	db := GetDatabase()
	result := db.Offset(offset_num).Order("follow_num desc").Limit(show_num).Find(&category_list)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}
	return
}

func CheckCategoryRepeat(name string) error {
	db := GetDatabase()
	result := db.Where("name = ?", name).Limit(1).Find(&models.Category{})
	if result.Error != nil {
		return EQueryFailed
	}
	if result.RowsAffected != 0 {
		return EHadExist
	}
	return nil
}

func CreateCategory(name string) (category *models.Category, err error) {
	if err = CheckCategoryRepeat(name); err != nil {
		return
	}
	category = &models.Category{
		Name:      name,
		FollowNum: 0,
	}

	db := GetDatabase()
	if db.Create(category).Error != nil {
		err = ECreateFailed
		return
	}
	return
}

func RenameCategory(categoryId uint, new_name string) error {
	db := GetDatabase()
	if db.Model(&models.Category{}).Update("name", new_name).Error != nil {
		return EUpdateFailed
	}
	return nil
}

func DeleteCategory(categoryId uint) error {
	db := GetDatabase()
	if db.Delete(&models.Category{}, categoryId).Error != nil {
		return EDeleteFailed
	}
	return nil
}
