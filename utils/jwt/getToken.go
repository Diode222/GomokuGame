package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(240)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte("webapp"))
	if err != nil {
		return ""
	}
	return tokenString
}
