package main

import (
	"captcha-backend/config"
	"captcha-backend/endpoint"
	"captcha-backend/service"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var conf config.Config
	if _, err := toml.DecodeFile("./config/config.toml", &conf); err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,
	})

	captchaService := service.NewCaptchaService()

	handler := endpoint.NewHandler(rdb, captchaService, &conf)

	err := fasthttp.ListenAndServe(":"+strconv.Itoa(conf.Port), handler.Endpoint)
	if err != nil {
		panic(err)
	}
}
