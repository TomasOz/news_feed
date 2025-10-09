package cache

import (
	"fmt"
)

func FeedKey(userID uint) string {
	return fmt.Sprintf("feed:%d", userID)
}


