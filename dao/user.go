package dao

import (
	"golang.org/x/crypto/bcrypt"
	"kome/mybbs-server/models"
)

func QueryUserbyId(userId uint) (user *models.User, err error) {
	db := GetDatabase()
	user = new(models.User)
	result := db.Where("id = ?", userId).Limit(1).Find(user)
	if result.Error != nil {
		err = EQueryFailed
		return
	}
	if result.RowsAffected == 0 {
		user.ID = userId
		user.Name = "UserNotFound"
		user.AgreeNum = 0
		user.StarNum = 0
		err = ENotExist
		return
	}
	return
}

func CheckUserNameRepeat(name string) error {
	db := GetDatabase()
	result := db.Where("name = ?", name).Limit(1).Find(&models.User{})
	if result.Error != nil {
		return EQueryFailed
	}
	if result.RowsAffected == 0 {
		return EHadExist
	}
	return nil
}

func CheckUserEmailRepeat(email string) error {
	db := GetDatabase()
	result := db.Where("email = ?", email).Limit(1).Find(&models.User{})
	if result.Error != nil {
		return EQueryFailed
	}
	if result.RowsAffected == 0 {
		return EHadExist
	}
	return nil
}

func CreateUser(name string, email string, password_origin string) (user *models.User, err error) {
	if len(password_origin) < 6 {
		err = EPasswordInvalid
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(password_origin), bcrypt.DefaultCost)
	if err != nil {
		err = EPasswordInvalid
		return
	}

	if err = CheckUserEmailRepeat(email); err != nil {
		return
	}
	if err = CheckUserNameRepeat(name); err != nil {
		return
	}

	user = &models.User{
		Name:     name,
		Email:    email,
		Password: string(password),
		AgreeNum: 0,
		StarNum:  0,
	}

	db := GetDatabase()
	if db.Create(user).Error != nil {
		err = ECreateFailed
		return
	}
	return
}

func DeleteUser(userId uint) error {
	db := GetDatabase()
	if db.Unscoped().Delete(&models.User{}, userId).Error != nil {
		return EDeleteFailed
	}
	return nil
}

func UpdateUser(userId uint, name string, email string, password_origin string, password_origin_old string) (user *models.User, err error) {
	db := GetDatabase()
	user, err = QueryUserbyId(userId)
	if err != nil {
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password_origin_old)); err != nil {
		err = EPasswordInvalid
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(password_origin), bcrypt.DefaultCost)
	if err != nil {
		err = EPasswordInvalid
		return
	}

	if email != user.Email && CheckUserEmailRepeat(email) != nil {
		err = EHadExist
		return
	}

	if name != user.Name && CheckUserNameRepeat(name) != nil {
		err = EHadExist
		return
	}

	user.Email = email
	user.Name = name
	user.Password = string(password)
	if err = db.Save(user).Error; err != nil {
		err = EUpdateFailed
		return
	}
	return
}

func UserLogin(email string, password_origin string) (user *models.User, err error) {
	db := GetDatabase()
	user = new(models.User)
	result := db.Where("email = ?", email).Limit(1).Find(user)
	if result.Error != nil {
		err = EQueryFailed
		return
	}

	if result.RowsAffected == 0 {
		err = ENotExist
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password_origin))
	if err != nil {
		err = EPasswordInvalid
		return
	}
	return
}
