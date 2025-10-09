package feed

import "news-feed/internal/post"

type FeedResponse struct {
    Posts    []post.Post `json:"posts"`
    HasMore  bool        `json:"has_more"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}