package follow

type UserFollows struct {
	// gorm.Model is skiped here, because it automatically adds 
    // ID        uint           `gorm:"primaryKey"`
    // CreatedAt time.Time
    // UpdatedAt time.Time
    // DeletedAt gorm.DeletedAt `gorm:"index"`
	// At this point it is not needed

    FolloweeID  uint `gorm:"primaryKey;not null;index:idx_followee" json:"followee_id"`
    FollowerID  uint `gorm:"primaryKey;not null;index:idx_follower" json:"follower_id"`

	// I skip FK association because foreign key checks add write overhead
	// User validation will be handled at code level
}
