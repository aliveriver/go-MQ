package main

import (
	"go-MQ/api"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	userapi := r.Group("/users")
	userapi.POST("/register", api.UserAPIEntity.RegisterHandler)
	userapi.POST("/login", api.UserAPIEntity.LoginHandler)
	return r
}
