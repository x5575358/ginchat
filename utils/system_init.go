package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *gorm.DB
	Rediss *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app init")
}

func InitRedis() {
	Rediss = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Rediss.Ping(Rediss.Context()).Result()
	if err != nil {
		fmt.Println("config redis fail", err)
	} else {
		fmt.Println("config redis init", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	//var err error
	err := Rediss.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Rediss.Subscribe(ctx, channel)
	fmt.Println("subscribe....",ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return msg.Payload, err
}

func InitMysql() {
	//自定义日志模板 打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	fmt.Println(newLogger)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	// fmt.Println(DB)
	fmt.Println("config mysql init")

}
