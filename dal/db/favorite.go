package db

import (
	"gorm.io/gorm"
)

// Favorite Gorm Data Structures
type Favorite struct {
	gorm.Model
	UserId  int64 `gorm:"column:user_id;not null;index:idx_userid"`
	VideoId int64 `gorm:"column:video_id;not null;index:idx_videoid"`
}

func (Favorite) TableName() string {
	return "favorite"
}
