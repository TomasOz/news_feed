package main

import (

	"news-feed/internal/user"
	"news-feed/internal/db"
	"github.com/gin-gonic/gin"

	"log"
)

func main() {
	router := gin.Default()

	dbConn := db.InitDatabase()

	userRepo := user.NewUserRepository(dbConn)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)

	log.Println("Lets Start Application")
	router.Run(":8080")

}
