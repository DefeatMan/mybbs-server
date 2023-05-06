package dao

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"kome/mybbs-server/models"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var kDB *gorm.DB

// Init Database, Create Datatable & Admin
func Init() {
	host := viper.Get("mysql.host")
	port := viper.Get("mysql.port")
	database := viper.Get("mysql.database")
	username := viper.Get("mysql.username")
	password := viper.Get("mysql.password")
	charset := viper.Get("mysql.charset")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, port, database, charset)
	log.Print("connect dsn is :", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Init Error: failed to connect mysql databse")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.UserStarPost{})
	db.AutoMigrate(&models.UserStarComment{})
	db.AutoMigrate(&models.UserAgreeComment{})

	kDB = db
	InitAdmin()
}

func GetDatabase() *gorm.DB {
	return kDB
}

func InitAdmin() {
	name, email, password_origin := getRoot()
	db := GetDatabase()
	db.Transaction(func(tx *gorm.DB) error {
		password, err := bcrypt.GenerateFromPassword([]byte(password_origin), bcrypt.DefaultCost)
		if err != nil {
			return EPasswordInvalid
		}

		init_user := models.User{
			Name: name, Email: email, Password: string(password),
			AgreeNum: 0, StarNum: 0,
		}
		if err := tx.Create(&init_user).Error; err != nil {
			return ECreateFailed
		}
		if err := tx.Create(&models.Admin{UserId: init_user.ID, AdminPerm: models.RootPermFlag}).Error; err != nil {
			return ECreateFailed
		}
		return nil
	})
}

func getRoot() (name string, email string, password string) {
	name = viper.Get("root.name").(string)
	email = viper.Get("root.email").(string)
	password = viper.Get("root.password").(string)
	return
}
