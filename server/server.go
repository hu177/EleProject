package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
)

var (
	rdb *redis.Client
)

const (
	RedisServer   = "127.0.0.1:6379"
	GinServerPort = ":8080"
)

func init() {
	//	初始化redis
	log.Println("start to init redis")
	rdb = redis.NewClient(&redis.Options{
		Addr:     RedisServer,
		Password: "",
		DB:       0,
	})
	ctx := context.TODO()
	if pong, err := rdb.Ping(ctx).Result(); err != nil || pong != "PONG" {
		log.Println("Redis init error : pong: ", err, pong)
		return
	}
}

func main() {
	// 开启gin端口，处理http请求
	api := gin.Default()
	// 发送接口，发送数据给小程序
	api.GET("/get", GetAxis)
	// 暂时监听本地端口
	err := api.Run(GinServerPort)
	log.Println("run error:", err)
}

func GetAxis(c *gin.Context) {
	// 从redis中取得key
	rKey, err := rdb.LPop(c, "TimeList").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"x": 44,
			"y": 120,
		})
		return
	}
	rdb.LPush(c, "TimeList", rKey)
	// 用key在maphash中取得坐标
	if rKey != "" {
		Ax, _ := rdb.HGet(c, rKey, "X").Result()
		Ay, _ := rdb.HGet(c, rKey, "Y").Result()
		XFloat, _ := strconv.ParseFloat(Ax, 64)
		YFloat, _ := strconv.ParseFloat(Ay, 64)
		c.JSON(http.StatusOK, gin.H{
			"x": XFloat,
			"y": YFloat,
		})
	}
	return
}
