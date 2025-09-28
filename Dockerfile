# ---------- Build stage ----------
    FROM golang:1.23 AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

    COPY . .
    
    RUN go build -o news-feed ./cmd
    
    FROM debian:bookworm-slim
    
    RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
    
    WORKDIR /app
    
    COPY --from=builder /app/news-feed .
    COPY --from=builder /app/migrations ./migrations
    COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
    EXPOSE 8080
    
    CMD ["./news-feed"]