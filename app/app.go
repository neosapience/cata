package main

import (
	"os"

	"cata/app/api"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func main() {
	route := setupRouter()
	route.Run(":8080")
}

func setupRouter() *gin.Engine {
	adminId := os.Getenv("ADMIN_ID")
	adminPass := os.Getenv("ADMIN_PW")
	redisAddr := os.Getenv("REDIS_ADDR")
	route := gin.Default()

	svc := &api.ApiBase{RedisCli: redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}

	authorized := route.Group("/api", gin.BasicAuth(gin.Accounts{
		adminId: adminPass,
	}))
	{
		authorized.GET("/cb", svc.GetListApiCb)
		authorized.POST("/cb", svc.PostApiCb)
		authorized.DELETE("/cb/:id", svc.DeleteApiCb)
		authorized.GET("/cb/:id", svc.GetOneApiCb)
		authorized.GET("/cbit", svc.GetListApiCbit)
	}

	anybody := route.Group("/api/cb")
	{
		anybody.POST("/:id", svc.PostOneApiCb)
	}

	return route
}
