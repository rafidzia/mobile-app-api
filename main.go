package main

import (
	"os"

	"github.com/alan1420/mobile-app-api/api/auth"
	"github.com/alan1420/mobile-app-api/api/user"
	"github.com/alan1420/mobile-app-api/constant"
	"github.com/alan1420/mobile-app-api/database"
	"github.com/alan1420/mobile-app-api/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// load .env begin
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// load .env end

	// database begin
	db, err := database.OpenDatabase()
	if err != nil {
		panic(err)
	}

	// repositories start
	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db, authRepository)
	// repositories end

	// routers start
	router := gin.Default()
	router.Use(CORSMiddleware())

	user.NewUserHandler(router, userRepository)
	auth.NewAuthHandler(router, authRepository)
	// routers end

	// run begin
	if os.Getenv("BASE_URL") != "" {
		constant.BaseUrl = os.Getenv("BASE_URL")
		constant.BaseUrlHttp = "http://" + constant.BaseUrl
	}
	err = router.Run(constant.BaseUrl)
	if err != nil {
		panic(err)
	}
	// run end
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
