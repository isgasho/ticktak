package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/api"
    "github.com/nihileon/ticktak/auth"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/middlewares"
)

func main() {
    log.Init()

    // config
    config, err := InitConfig()
    if err != nil {
        panic(fmt.Errorf("init config error: %s", err))
    }

    // MySQL
    err = dal.InitDB(config.MysqlDSN)
    if err != nil {
        panic(err)
    }

    // Redis Or Map
    err = dal.InitKV(config.RedisAddr, config.MemoryOrRedis)
    if err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(middlewares.AllowControl())

    r.POST("/register", auth.RegisterUser)
    r.POST("/login", auth.Login)

    r.Use(middlewares.JWTAuth())

    r.OPTIONS("api/*pattern", middlewares.OPTIONSHandle)

    platformAPI := r.Group("api")
    platformAPI.POST("/user/update", api.ChangeCurrentUser)
    platformAPI.GET("/task/list", api.GetTasksByUsername)
    platformAPI.POST("/task/update/state", api.ChangeTaskState)
    platformAPI.POST("/task/update/priority", api.ChangeTaskPriority)
    platformAPI.POST("/task/add", api.AddTask)
    platformAPI.GET("/task/list/state/", api.GetTasksByUsernameState)
    platformAPI.GET("/task/list/priority/", api.GetTasksByUsernamePriority)

    r.Run(config.ListenAddr)
}
