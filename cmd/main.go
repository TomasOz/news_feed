// @title           News Feed API
// @version         1.0
// @description     News feed API (learning project).
// @contact.name    Tomas Ozolinsius
// @contact.url     https://github.com/TomasOz
// @BasePath        /

// Security: JWT via Authorization: Bearer <token>
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"news-feed/cmd/docs"

	"news-feed/internal/user"
	"news-feed/internal/post"
	"news-feed/internal/follow"
	"news-feed/internal/feed"

	"news-feed/internal/db"
	"news-feed/internal/auth"
	"news-feed/internal/cache"
	"news-feed/internal/background_jobs"

	"news-feed/internal/health"
)

func main() {
	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()

	dbConn := db.InitDatabase()

	// Initialize Redis cache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	
	cacheClient := cache.NewRedisCache(redisAddr, "", 0)
	log.Println("Connected to Redis cache")

	// REPOS
	userRepo := user.NewUserRepository(dbConn)
	postRepo := post.NewPostRepository(dbConn)
	followRepo := follow.NewFollowRepository(dbConn)

	// SERVICES
	userService := user.NewUserService(userRepo)
	postService := post.NewPostService(postRepo)
	followService := follow.NewFollowService(followRepo)
	feedService := feed.NewFeedService(followRepo, postRepo, cacheClient)

	// BACKGROUND JOBS ==========================
	fanoutWorker := background_jobs.NewFanoutWorker(followRepo, cacheClient)
	fanoutWorker.Start()

	// HANDLERS
	userHandler := user.NewUserHandler(userService)
	postHandler := post.NewPostHandler(postService, fanoutWorker)
	followHandler := follow.NewFollowHandler(followService)
	feedHandler := feed.NewFeedHandler(feedService)


	//ROUTES
	//USER
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)

	// POST =====================================

	router.GET("/post/:id", auth.AuthMiddleware(), postHandler.GetPostByID)
	//Temporary for receiving all posts
	router.GET("/posts", auth.AuthMiddleware(), postHandler.GetPosts)
	router.POST("/posts", auth.AuthMiddleware(), postHandler.CreatePost)

	// FOLLOW =====================================
	// Users that following current user
	router.POST("/follow/:id", auth.AuthMiddleware(), followHandler.Follow)
	// What current user is following
	router.POST("/unfollow/:id", auth.AuthMiddleware(), followHandler.UnFollow)

	// FEED =====================================
	router.GET("/feed", auth.AuthMiddleware(), feedHandler.GetFeed)

	// SWAGGER ===================================
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", health.Health)

	router.GET("/ready", health.Readiness(cacheClient, dbConn))

	log.Println("Lets Start Application with Redis Caching and Background Jobs")
	router.Run(":8080")

	defer fanoutWorker.Stop()
}