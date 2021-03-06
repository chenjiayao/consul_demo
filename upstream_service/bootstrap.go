package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upstream_service/globals"
	"upstream_service/routes"
	"upstream_service/utils/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func run() {

	logger.InitLogger()
	initCache()
	// initDB()
	initHttpEngine()
	routes.LoadRoute()

	go func() {
		port := viper.GetInt("http.port")
		address := viper.GetString("http.address")
		globals.HttpServer = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", address, port),
			Handler: globals.Engine,
		}

		logger.Logger.Sugar().Infof("http server listen on %s:%d", address, port)
		err := globals.HttpServer.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()
	waitSignal()
}

func initHttpEngine() {
	globals.Engine = gin.Default()

	if !viper.GetBool("degbug") {
		gin.SetMode(gin.ReleaseMode)
	}

	globals.Engine = gin.New()
	globals.Engine.Use(gin.Recovery())

	if viper.GetBool("http.gzip_enable") {
		globals.Engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	if viper.GetBool("http.cors_enable") {
		globals.Engine.Use(cors.New(cors.Config{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
}

func initCache() {
	address := viper.GetString("redis.address")
	port := viper.GetInt32("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	globals.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Password: password,
		DB:       db,
	})

	for {
		if err := globals.RedisClient.Ping(context.TODO()).Err(); err != nil {
			logger.Logger.Sugar().Errorf("redis connect error: %s, retry after 30s", err)
			time.Sleep(time.Second * 30)
		} else {
			logger.Logger.Sugar().Infof("redis connect success")
			break
		}
	}
}

func initDB() {
	db, err := gorm.Open(mysql.New(mysql.Config{

		DSN:                       viper.GetString("db.dsn"),
		DefaultStringSize:         256,   // string ???????????????????????????
		DisableDatetimePrecision:  true,  // ?????? datetime ?????????MySQL 5.6 ???????????????????????????
		DontSupportRenameIndex:    true,  // ???????????????????????????????????????????????????MySQL 5.7 ????????????????????? MariaDB ????????????????????????
		DontSupportRenameColumn:   true,  // ??? `change` ???????????????MySQL 8 ????????????????????? MariaDB ?????????????????????
		SkipInitializeWithVersion: false, // ???????????? MySQL ??????????????????
	}), &gorm.Config{})
	if err != nil {
		logger.Logger.Sugar().Fatalf("db connect error: %s", err)
		os.Exit(1)
	}
	globals.DB = db
}

func waitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	err := globals.HttpServer.Shutdown(ctx)
	logger.Logger.Sugar().Infof("http server shutdown: %s", err)
	logger.Logger.Sync()

	if err != nil {
	}
}
