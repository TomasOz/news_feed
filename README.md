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
- **Cache:** Redis
- **Background Jobs:** Go goroutines + channels
- **Database Migrations:** golang-migrate
- **Documentation:** Swagger
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

### Run with Dockergit

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

## Docs (Swagger)

Project already delivers generated docs, but in case of regenerating doc run
```bash
swag init -g cmd/main.go -o cmd/docs --parseDependency --parseInternal
```

Docs available at 
http://localhost:8080/swagger/index.html


## Future Plans
- Add Redis caching (done)
- Introduce background workers (done)
- Add Swagger documentation (done)
- Add Testing

- Maybe lets think about solution, where feed would be deleted if it is read
- At the moment /feed uses just simple limit=10&offset=10 in real world application
it would best to have cursor pagination
- There should be followers cache also
- Finish Swagger documentation

## Author

Built by Tomas Ozolinsius