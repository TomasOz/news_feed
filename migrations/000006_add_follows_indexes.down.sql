-- Drop the indexes added in the up migration
DROP INDEX IF EXISTS idx_follows_follower ON user_follows;
