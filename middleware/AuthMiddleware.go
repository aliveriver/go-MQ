package middleware

import (
	"context"
	"errors"
	"go-MQ/common"
	"go-MQ/entity"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkToken(ctx context.Context, tokenString string) (*common.Claims, error) {
	//验证格式
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("未找到token或token格式错误")
	}
	tokenString = tokenString[7:]

	token, claim, err := common.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return nil, errors.New("token无效")
	}
	n, _ := common.GetRedisClient().Exists(ctx, "blacklist:"+claim.Id).Result()
	if n > 0 {
		return nil, errors.New("token已失效")
	}
	return claim, nil
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		claim, err := checkToken(ctx, tokenString)
		if err != nil {
			common.SendErrorResponse(ctx, err.Error())
			ctx.Abort()
			return
		}

		//验证通过，获取Token中的userid
		userId := claim.UserID
		DB := common.GetDB()
		var user entity.User
		DB.First(&user, userId)

		//判断用户存在
		if user.ID == 0 {
			common.SendErrorResponse(ctx, "用户不存在")
			ctx.Abort()
			return
		}

		//写入user信息
		ctx.Set("user", user)
		// 将 ExpiresAt (秒) 转为毫秒（13位），若为0则写0
		var expiresAtMs int64
		if claim != nil && claim.ExpiresAt > 0 {
			expiresAtMs = claim.ExpiresAt * 1000
		} else {
			expiresAtMs = 0
		}
		ctx.Set("TokenExpiresAt", expiresAtMs)
		ctx.Set("jti", claim.Id)

		ctx.Next()
	}
}
