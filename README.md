# News Feed Server

## Overview

A simple News feed API written in Go:
This is learning system architecture learning project written in Golang. 
I will use a stack that I’m not familiar with in my daily work.

### Core features:

- User registration and login with JWT authentication
- Post creation (news updates)
- Follow/unfollow other users
- Fetch news feed from following users

## Used Stack
- **Language:** Golang
- **Web framework:** Gin
- **ORM:** GORM  
- **Database:** MySQL
- **Auth:** JWT  
- **Containerization:** Docker & Docker Compose

## How to Start

### Locally:

 Required

- Go 1.23+
- MySQL server local or remote 

```bash
    go run cmd/main.go
```

### Run with Docker

```bash
    docker compose up 
```

### Environment Variables

```bash
    MYSQL_ROOT_PASSWORD=yourpassword
    MYSQL_DATABASE=newsfeed
    DB_URL=root:changeme@tcp(db:3306)/newsfeed?parseTime=true
    JWT_SECRET=your_super_secret_key
```

App will be available at
http://localhost:8080


## Future Plans
- Add Redis caching
- Introduce background workers
- Add Swagger documentation
- Add Testing
- I will add more later 

## Author

Built by Tomas Ozolinsius