package common

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	return tokenString, expirationTime.UnixMilli(), nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}

func AddInvalidToken(jti string, ExpiresAt int64) error {
	// 2. 计算剩余有效期 (ExpiresAt - Now)
	expiration := time.Unix(ExpiresAt, 0).Sub(time.Now())

	// 如果已经过期，无需加入黑名单
	if expiration <= 0 {
		return nil
	}

	// 3. 存入 Redis
	// Key: "blacklist:eyJxh..." (建议加前缀)
	// Value: "1" (值无所谓)
	// Expiration: 动态计算出的剩余时间
	err := GetRedisClient().Set(context.Background(), "blacklist:"+jti, "1", expiration).Err()
	return err
}
