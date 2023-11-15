package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis服务器地址
		Password: "",               // Redis密码，如果没有设置密码，可以为空字符串
		DB:       0,                // Redis数据库索引
	})

	// 测试Redis连接
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("连接Redis失败:", err)
		return
	}
	fmt.Println("Redis连接成功:", pong)
	RedisClient = client

}

func GetStatus() (int, error) {
	// 读取指定的值
	value, err := RedisClient.Get("qskm-backend:status").Int()
	if err != nil {
		fmt.Println("读取Redis值失败:", err)
		return 0, err
	}
	return value, nil
}
