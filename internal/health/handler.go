package health

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"context"
	"gorm.io/gorm"

	"news-feed/internal/cache"

	"database/sql"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "Service is running",
	})
}
	
func Readiness(redis *cache.RedisCache, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		err := redis.Ping(ctx)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"message": "Redis connection failed",
				"error": err.Error(),
			})
			return
		}

		var dbErr error
		if db != nil {
			var sqlDB *sql.DB
			if sqlDB, dbErr = db.DB(); dbErr == nil {
				dbErr = sqlDB.PingContext(ctx)
			}
		}

		if dbErr != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"message": "Mysql connection failed",
				"error": dbErr.Error(),
			})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"message": "All Services are running",
		})
	}
}

