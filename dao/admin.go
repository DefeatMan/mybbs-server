package dao

import (
	"kome/mybbs-server/models"
)

func QueryAdminbyUserId(userId uint) (admin *models.Admin, err error) {
	db := GetDatabase()
	admin = new(models.Admin)
	result := db.Where("user_id = ?", userId).Limit(1).Find(admin)
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

func AppendAdminPerm(userId uint, adminPerm uint) error {
	db := GetDatabase()
	result, err := QueryAdminbyUserId(userId)
	if err != nil {
		return err
	}
	result.AdminPerm = result.AdminPerm | adminPerm
	if err = db.Model(result).Update("admin_perm", result.AdminPerm).Error; err != nil {
		return EUpdateFailed
	}
	return nil
}

func SubsetAdminPerm(userId uint, adminPerm uint) error {
	db := GetDatabase()
	result, err := QueryAdminbyUserId(userId)
	if err != nil {
		return err
	}
	result.AdminPerm = result.AdminPerm & adminPerm
	if err = db.Model(&result).Update("admin_perm", result.AdminPerm).Error; err != nil {
		return EUpdateFailed
	}
	return nil
}

func CreateAdmin(userId uint, adminPerm uint) error {
	db := GetDatabase()
	var had int64 = 0
	err := db.Where("user_id = ?", userId).Limit(1).Count(&had)
	if err.Error != nil {
		return EQueryFailed
	}
	if had == 0 {
		return ENotExist
	}
	create_info := models.Admin{
		UserId:    userId,
		AdminPerm: adminPerm,
	}
	if err = db.Create(&create_info); err.Error != nil {
		return ECreateFailed
	}
	return nil
}

func DeleteAdmin(userId uint) error {
	db := GetDatabase()
	err := db.Unscoped().Where("user_id = ?", userId).Delete(&models.Admin{})
	if err.Error != nil {
		return EDeleteFailed
	}
	if err.RowsAffected == 0 {
		return ENotExist
	}
	return nil
}
