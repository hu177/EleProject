package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	local       = "127.0.0.1:8888"
	server      = ":8888"
	RedisServer = "127.0.0.1:6379"
	RedisTime   = "2006-01-02 15:04:05"
)

var (
	rdb *redis.Client
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
//test 
	fmt.Println("监听端口发来的消息")
	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	// 循环开启监听端口，接收tcp请求 --> 这里指位置信息
	for {
		conn, _ := listener.Accept()
		go handleConn(conn)
	}
}

// 接收数据存到转存到服务器redis上
func handleConn(conn net.Conn) {
	//循环不停的去处理数据
	for {
		// 使用字符切片去接收数据
		tmpByte := make([]byte, 1)
		var resData []byte
		//循环去读取数据
		for {
			length, _ := conn.Read(tmpByte)
			//fmt.Println(length, tmpByte
			//读到的长度为0，或者读取到换行就断掉
			if length == 0 || tmpByte[0] == 10 {
				break
			}
			//拼接读到的数据,
			resData = append(resData, tmpByte...)

		}
		str := fmt.Sprintf("%s", string(resData))

		// str即为接收到的坐标字符串
		if len(str) != 0 {
			axises := strings.Split(str, ",")
			if len(axises) != 2 {
				log.Println("axises get error:", str)
				return
			}
			Xaxis := axises[0][2:]
			Yaxis := axises[1][2:]
			XFloat, _ := strconv.ParseFloat(Xaxis, 64)
			YFloat, _ := strconv.ParseFloat(Yaxis, 64)
			fmt.Println(XFloat, YFloat)
			//	存入redis，当前时间为key
			ctx := context.TODO()
			keyTime := time.Now().Format(RedisTime)
			// 使用
			rdb.LPush(ctx, "TimeList", keyTime)
			// 存入hashkey中
			err := rdb.HSet(ctx, keyTime, "X", XFloat, "Y", YFloat).Err()
			if err != nil {
				log.Println("[GeoAdd error:]", err)
			}
		}
		//fmt.Println("str:\n",str)
	}
}
