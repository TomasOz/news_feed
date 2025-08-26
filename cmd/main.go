package main

import (

	"news-feed/internal/user"
	"news-feed/internal/post"
	"news-feed/internal/db"
	"news-feed/internal/auth"
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

	postRepo := post.NewPostRepository(dbConn)
	postService := post.NewPostService(postRepo)
	postHandler := post.NewPostHandler(postService)

	router.GET("/post/:id", auth.AuthMiddleware(), postHandler.GetPostByID)
	//Temporary for receiving all posts
	router.GET("/posts", auth.AuthMiddleware(), postHandler.GetPosts)
	router.POST("/posts", auth.AuthMiddleware(), postHandler.CreatePost)

	log.Println("Lets Start Application")
	router.Run(":8080")

}
