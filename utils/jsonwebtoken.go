package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Claims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// default secret string
const dSecret = "alwayswinwinwin"

var kSecret []byte
var kExpireDuration = time.Hour * 2

func InitJwt() (err error) {
	vSecret := viper.Get("token.secret")
	if vSecret == nil {
		fmt.Println("[init jwt warn] token.secret nil")
        kSecret = []byte(dSecret)
	} else {
        kSecret = []byte(vSecret.(string))
    }
	fmt.Println("kSecret: ", kSecret)

	vExpire := viper.Get("token.expire")
	if vExpire == nil {
		fmt.Println("[init jwt warn] token.expire nil")
	} else {
		kExpireDuration = time.Duration(int(vExpire.(float64))) * time.Minute
	}
	fmt.Println("kExpertDuration: ", kExpireDuration)
	return
}

func GetJwtKey(_ *jwt.Token) (i interface{}, err error) {
	return kSecret, nil
}

func NewToken(userId uint, email string) (tokenS string, err error) {
	c := Claims{
		UserId: userId,
		Email:  email,
        StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(kExpireDuration).Unix(), // expiration time
			Issuer:    "myserver",
		},
	}

	// create new token with signed secret
	tokenS, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(kSecret)
	return
}

func ParseToken(tokenS string) (claims *Claims, err error) {
	var token *jwt.Token
	claims = new(Claims)
	token, err = jwt.ParseWithClaims(tokenS, claims, GetJwtKey)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

func RefreshToken(tokenS string) (newTokenS string, err error) {
	var claims Claims
	_, err = jwt.ParseWithClaims(tokenS, claims, GetJwtKey)
	v, _ := err.(*jwt.ValidationError)
    
	// create new token at expiration
	if v.Errors == jwt.ValidationErrorExpired {
		return NewToken(claims.UserId, claims.Email)
	}
	return
}
