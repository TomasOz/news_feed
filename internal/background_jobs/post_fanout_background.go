package background_jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"news-feed/internal/cache"
	"news-feed/internal/follow"
)

type FanoutJob struct {
	PostID     uint      `json:"post_id"`
	FollowerID uint      `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FanoutWorker struct {
	followRepo follow.FollowRepository
	cache      *cache.RedisCache
	jobQueue   chan FanoutJob
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewFanoutWorker(followRepo follow.FollowRepository, cache *cache.RedisCache) *FanoutWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &FanoutWorker{
		followRepo: followRepo,
		cache:      cache,
		jobQueue:   make(chan FanoutJob, 1000), // Buffer for 1000 jobs
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (w *FanoutWorker) Start() {
	log.Println("Starting fanout worker...!!!!!!!!!!!!!!!!!!!")
	
	// Start multiple workers for parallel processing
	for i := 0; i < 5; i++ {
		go w.worker(i)
	}
}

func (w *FanoutWorker) Stop() {
	log.Println("Stopping fanout worker...")
	w.cancel()
	close(w.jobQueue)
}

func (w *FanoutWorker) worker(workerID int) {
	log.Printf("Fanout worker %d started", workerID)
	
	for {
		select {
		case job := <-w.jobQueue:
			if err := w.processFanout(job); err != nil {
				log.Printf("Worker %d failed to process job: %v", workerID, err)
			}
		case <-w.ctx.Done():
			log.Printf("Fanout worker %d stopped", workerID)
			return
		}
	}
}

func (w *FanoutWorker) QueueFanoutJob(job FanoutJob) {
	select {
	case w.jobQueue <- job:
		log.Printf("Queued fanout job for follower %d, post %d", job.FollowerID, job.PostID)
	default:
		log.Printf("Job queue full, dropping fanout job for follower %d", job.FollowerID)
	}
}

func (w *FanoutWorker) processFanout(job FanoutJob) error {
	ctx := context.Background()
	
	feedKey := cache.FeedKey(job.FollowerID)
	postIDStr := fmt.Sprintf("%d", job.PostID)
	
	// Lets push value to head
	// TO DO in the future there should be some kind of expiration time
	if err := w.cache.LPush(ctx, feedKey, postIDStr); err != nil {
		return fmt.Errorf("failed to add post to feed cache: %w", err)
	}
	
	// Lets Keep 1000 values only
	if err := w.cache.LTrim(ctx, feedKey, 0, 999); err != nil {
		log.Printf("Warning: failed to trim feed cache for user %d: %v", job.FollowerID, err)
	}
		
	return nil
}

func (w *FanoutWorker) FanoutToAllFollowers(postID uint, authorID uint) error {
	followers, err := w.followRepo.GetFollowersID(authorID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}
	
	if len(followers) == 0 {
		log.Printf("No followers found for user %d", authorID)
		return nil
	}
	
	// Create a job for each follower
	for _, followerID := range followers {
		job := FanoutJob{
			PostID:     postID,
			FollowerID: followerID,
			CreatedAt:  time.Now(),
		}
		
		w.QueueFanoutJob(job)
	}
	
	log.Printf("Created %d fanout jobs for post %d", len(followers), postID)
	return nil
}