package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	_ "gin-reaper/docs"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.NewApp()

	// base application info
	app.Name = "reaper server"
	app.Author = "HBX"
	app.Version = "0.0.1"
	app.Copyright = "HBX team reserved"
	app.Usage = "reaper server"

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./configs/config.toml",
			Usage: "load configuration from `FILE`",
		},
	}

	app.Action = server
	app.Run(os.Args)
}

func server(c *cli.Context) error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(c.String("config"))
	viper.ReadInConfig()

	log.SetFormatter(&log.TextFormatter{DisableColors: true})

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() { defer wg.Done(); Startup() }()

	wg.Wait()
	return nil
}

func Startup() {
	rand.Seed(time.Now().Unix())

	// start gin
	route := gin.Default()

	route.GET("/test", Test)

	// swagger访问地址 http://localhost:8080/swagger/index.html
	//url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := route.Group("v1")
	{
		v1.GET("/login", login)
		v1.GET("/user", user)
	}

	v2 := route.Group("v2")
	{
		v2.GET("/login", login)
		v2.GET("/user", user)
	}

	route.Run(viper.GetString("webserver.addr"))

}

// @Summary login
// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /v1
// @Host 127.0.0.1:8080
// @Produce  json
// @Param name query string true "Name"
// @Success 200 {string} json "{"code":200,"data":"name","msg":"ok"}"
// @Router /v1 [get]
func login(ctx *gin.Context) {
	fmt.Println("this is a login method!")
}

// @Summary user api
// @Description this is user api
// @Success 200 {string} json "{"name": string, "message": string}"
// @Router /v1/user [get]
func user(ctx *gin.Context) {
	fmt.Println("this is a user method!")
	ctx.JSON(http.StatusOK, gin.H{
		"name":    "zbb",
		"message": "hello",
	})
}

// @Summary 测试接口
// @Description 描述信息
// @Success 200 {string} string  "ok"
// @Router /test [get]
func Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
