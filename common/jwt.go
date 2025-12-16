package common

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var jwtKey = []byte(viper.GetString("jwt.jwtkey"))

type Claims struct {
	UserID uint64
	jwt.StandardClaims
}

func ReleaseToken(userID uint64) (string, int64, error) {
	expirationTime := time.Now().Add(3 * 24 * time.Hour) //token有效期为3天
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), //token发放时间
			Issuer:    "dkd",             //谁发放的token
			Subject:   "user token",      //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}
	key := "usertoken:" + strconv.FormatUint(claims.UserID, 10)
	err = GetRedisClient().Set(context.Background(), key, claims.Id, time.Until(expirationTime)).Err()
	if err != nil {
		logrus.Errorf("redis set token error: %v", err)
		return "", 0, err
	}
	return tokenString, expirationTime.UnixMilli(), nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}

func AddInvalidToken(userID uint64) error {
	key := "usertoken:" + strconv.FormatUint(userID, 10)

	// 直接删除 key
	err := GetRedisClient().Del(context.Background(), key).Err()
	if err != nil {
		// 处理错误
		logrus.Errorf("redis del token error: %v", err)
		return err
	}
	return nil
}
