package follow

import (
	"errors"
)

type FollowService interface {
	Follow(follower_id, followee_id uint) (error)
	UnFollow(follower_id, followee_id uint) (error)
}

type DefaultFollowService struct {
	repo FollowRepository
}

func NewFollowService(repo FollowRepository) FollowService {
	return &DefaultFollowService{repo: repo}
}

func (p DefaultFollowService) Follow(follower_id, followee_id uint) (error) {
	if follower_id == followee_id {
		return errors.New("you can not follow yourself")
	}

	alreadyFollowing, err := p.repo.AlreadyFollowing(follower_id, followee_id)
	
	if err != nil {
		return err
	}

	if alreadyFollowing {
		return errors.New("you already follow user")
	}
	
	err = p.repo.Follow(follower_id, followee_id)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultFollowService) UnFollow(follower_id, followee_id uint) (error) {
	if follower_id == followee_id {
		return errors.New("you can not unfollow yourself")
	}
	
	err := p.repo.UnFollow(follower_id, followee_id)

	if err != nil {
		return errors.New("user was not found")
	}

	return nil
}