package main

import (

	"news-feed/internal/user"
	"news-feed/internal/post"
	"news-feed/internal/follow"
	"news-feed/internal/feed"

	"news-feed/internal/db"
	"news-feed/internal/auth"
	"github.com/gin-gonic/gin"

	"log"
)

func main() {
	router := gin.Default()

	dbConn := db.InitDatabase()

	// USER =====================================
	userRepo := user.NewUserRepository(dbConn)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)

	// POST =====================================
	postRepo := post.NewPostRepository(dbConn)
	postService := post.NewPostService(postRepo)
	postHandler := post.NewPostHandler(postService)

	router.GET("/post/:id", auth.AuthMiddleware(), postHandler.GetPostByID)
	//Temporary for receiving all posts
	router.GET("/posts", auth.AuthMiddleware(), postHandler.GetPosts)
	router.POST("/posts", auth.AuthMiddleware(), postHandler.CreatePost)

	// FOLLOW =====================================
	followRepo := follow.NewFollowRepository(dbConn)
	followService := follow.NewFollowService(followRepo)
	followHandler := follow.NewFollowHandler(followService)

	// Users that following current user
	router.POST("/follow/:id", auth.AuthMiddleware(), followHandler.Follow)
	// What current user is following
	router.POST("/unfollow/:id", auth.AuthMiddleware(), followHandler.UnFollow)

	// FEED =====================================
	feedService := feed.NewFeedService(followRepo, postRepo)
	feedHandler := feed.NewFeedHandler(feedService)

	router.GET("/feed", auth.AuthMiddleware(), feedHandler.GetFeed)

	log.Println("Lets Start Application")
	router.Run(":8080")

}
