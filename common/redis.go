package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       0,
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		logrus.Error("redis connect failed")
	}
}

func GetRedisClient() *redis.Client {
	return rdb
}
