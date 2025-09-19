package follow

import (
	"gorm.io/gorm"
	"errors"
)

type FollowRepository interface {
	Follow(follower_id, followee_id uint) (error)
	UnFollow(follower_id, followee_id uint) (error)
	AlreadyFollowing(follower_id, followee_id uint) (bool)
	GetFolloweesID(follower_id uint) ([]uint, error)
}

type GormFollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &GormFollowRepository{db: db}
}

func (r *GormFollowRepository) Follow(follower_id, followee_id uint) (error) {
	userFollows := UserFollows{
		FollowerID: follower_id,
		FolloweeID: followee_id,
	}

	if err := r.db.Create(&userFollows).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormFollowRepository) UnFollow(follower_id, followee_id uint) (error) {
	if err := r.db.Where("follower_id = ? AND followee_id = ?", follower_id, followee_id).Delete(&UserFollows{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormFollowRepository) AlreadyFollowing(follower_id, followee_id uint) (bool) {
	if err := r.db.Where("follower_id = ? AND followee_id = ?", follower_id, followee_id).First(&UserFollows{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}

	}
	return true
}

func (r *GormFollowRepository) GetFolloweesID(follower_id uint) ([]uint, error) {
	var falloweesIds []uint
	
	err := r.db.Model(&UserFollows{}).
		Where("follower_id = ? ", follower_id).
		Pluck("followee_id", &falloweesIds).Error

	if err != nil {
		return nil, err
	} 

	
	return falloweesIds, nil
}