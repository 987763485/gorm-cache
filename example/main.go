package main

import (
	"context"
	"fmt"
	"github.com/987763485/gorm-cache/store/gormredis"
	"github.com/987763485/gorm-cache/v1"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	db *gorm.DB

	redisClient *redis.Client

	cachePlugin *gormcache.Cache
)

func newDb() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	redisClient = redis.NewClient(&redis.Options{Addr: ":6379"})

	cacheConfig := &gormcache.Config{
		Store:      gormredis.NewWithDb(redisClient), // OR redis2.New(&redis.Options{Addr:"6379"})
		Serializer: &gormcache.DefaultJSONSerializer{},
	}

	cachePlugin = gormcache.New(cacheConfig)

	if err = db.Use(cachePlugin); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func basic() {
	var username string
	ctx := context.Background()
	ctx = gormcache.NewExpiration(ctx, time.Hour)

	db.Table("users").WithContext(ctx).Where("id = 1").Limit(1).Pluck("username", &username)
	fmt.Println(username)
	// output gorm
}

func customKey() {
	var nickname string
	ctx := context.Background()
	ctx = gormcache.NewExpiration(ctx, time.Hour)
	ctx = gormcache.NewKey(ctx, "nickname")

	db.Table("users").WithContext(ctx).Where("id = 1").Limit(1).Pluck("nickname", &nickname)

	fmt.Println(nickname)
	// output gormwithmysql
}

func useTag() {
	var nickname string
	ctx := context.Background()
	ctx = gormcache.NewExpiration(ctx, time.Hour)
	ctx = gormcache.NewTag(ctx, "users")

	db.Table("users").WithContext(ctx).Where("id = 1").Limit(1).Pluck("nickname", &nickname)

	fmt.Println(nickname)
	// output gormwithmysql
}

func main() {
	newDb()
	basic()
	customKey()
	useTag()

	ctx := context.Background()
	fmt.Println(redisClient.Keys(ctx, "*").Val())

	fmt.Println(cachePlugin.RemoveFromTag(ctx, "users"))
}
