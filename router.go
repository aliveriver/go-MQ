package main

import (
	"go-MQ/api"
	"go-MQ/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	userapi := r.Group("/users")
	userapi.POST("/register", api.UserAPIEntity.Register)
	userapi.POST("/login", api.UserAPIEntity.Login)
	userapi.POST("/logout", middleware.UserAuthMiddleware(), api.UserAPIEntity.Logout)
	userapi.PUT("/update", middleware.UserAuthMiddleware(), api.UserAPIEntity.UpdateUserInfo)
	return r
}
