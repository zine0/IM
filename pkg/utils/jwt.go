package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	Uid      int `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, uid int,secretKey string) (string, error) {
	
	claims := UserClaims{
		username,
		uid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return t.SignedString([]byte(secretKey))
}

func ValidJWT(j string,secretKey string) (*UserClaims,error) {

	t,err := jwt.ParseWithClaims(j,&UserClaims{},func(t *jwt.Token) (any, error) {
		return []byte(secretKey),nil
	})

	if claims,ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims,nil
	}else {
		return nil,err
	}
}
