package db

import (
    "log"
    "os"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "news-feed/migrations/migration"
)

func InitDatabase() *gorm.DB {
    dsn := os.Getenv("DB_URL")

    var db *gorm.DB
    var err error

	for i := 1; i <= 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect MySQL after retries:", err)
	}

    // if err = db.AutoMigrate(&user.User{}, &post.Post{}, &follow.UserFollows{}); err != nil {
    //     log.Fatal("Failed to auto-migrate:", err)
    // }

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get sql.DB:", err)
    }

    if err := migration.RunMigrations(sqlDB, "./migrations"); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    return db
}