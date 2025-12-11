package main

import (
	"bytes"
	"go-MQ/common"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {
	Init()

	r := gin.New()
	r.Use(gin.Recovery()) // default error handler
	r = CollectRoute(r)

	// enable swagger in debug mode
	if os.Getenv("SWAGGER_ENABLE") == "true" {
		docs.SwaggerInfo.Host = viper.GetString("server.host")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		logrus.Warn("swagger enabled, ensure you are in a safe environment:", "http://localhost:8080/swagger/index.html")
	}

	logrus.Info("server start")
	panic(r.Run(":8080"))
}

func Init() {
	InitConfig()
	common.InitLogrus()
	common.InitDB()
}

var configFile []byte

func InitConfig() {
	viper.SetConfigType("yml")
	// If a config file was embedded into the binary (configFile != nil), use it.
	// Otherwise attempt to read from `config/config.yml` on disk.
	if len(configFile) > 0 {
		err := viper.ReadConfig(bytes.NewBuffer(configFile))
		if err != nil {
			panic(err)
		}
		return
	}
	viper.SetConfigFile("config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
