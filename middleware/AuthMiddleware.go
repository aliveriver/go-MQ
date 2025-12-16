package middleware

import (
	"context"
	"errors"
	"go-MQ/common"
	"go-MQ/entity"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	// 获取 Redis 中的值
	key := "usertoken:" + strconv.FormatUint(claim.UserID, 10)
	val, err := common.GetRedisClient().Get(ctx, key).Result()

	if err == redis.Nil {
		// 1. Redis 中不存在该 Key -> 说明登录已过期或从未登录
		return nil, errors.New("token已失效或过期")
	} else if err != nil {
		// 2. Redis 连接或其他错误
		return nil, err
	}
	if val != claim.Id {
		return nil, errors.New("token已失效或过期")
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
