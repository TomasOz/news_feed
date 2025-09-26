package feed

import "news-feed/internal/post"

type FeedResponse struct {
    Posts      []post.Post 		   `json:"posts"`
    NextCursor string              `json:"next_cursor"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}