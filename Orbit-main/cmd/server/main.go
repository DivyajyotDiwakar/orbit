package main

import (
	"Orbit/configs"
	"Orbit/internal/db"
	routergroup "Orbit/internal/routerGroup"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	r := gin.Default()
	cfg := configs.LoadConfig()
	_ = db.GetInstance()
	_ = db.GetR2Client()
	_ = db.GetRedisClient()

	ApiGroup := r.Group("/api") // this is the route to group all the backend endpoints.
	routergroup.ApiRoutes(ApiGroup)

	server := &http.Server{
		Addr:    ":" + cfg.PORT.Value,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	GraceFulShutDown(server)
}
