CREATE INDEX idx_posts_feed_lookup ON posts (user_id, created_at DESC, id DESC);
