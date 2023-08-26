package db

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	UserId        int64  `gorm:"column:user_id;not null;index:idx_userid"` // 视频作者信息,绑定用户id，索引
	Title         string `gorm:"column:title;type:varchar(128);not null"`  // 视频标题
	PlayUrl       string `gorm:"column:play_url;varchar(128);not null"`    // 视频播放地址
	CoverUrl      string `gorm:"column:cover_url;varchar(128);not null"`   // 视频封面地址
	FavoriteCount int64  `gorm:"column:favorite_count;default:0"`          // 视频的点赞总数
	CommentCount  int64  `gorm:"column:comment_count;default:0"`           // 视频的评论总数
}

func (v *Video) TableName() string {
	return "video"
}

// QueryVideoByLatestTime query video info by latest create Time
func QueryVideoByLatestTime(ctx context.Context, latestTime int64) ([]*Video, error) {
	var videos []*Video
	updated_time := time.UnixMilli(latestTime)
	if err := DB.WithContext(ctx).Limit(30).Order("updated_at desc").Where(`updated_at < ?`, updated_time).Find(&videos).Error; err != nil {
		klog.Error("QueryVideoByLatestTime find video error " + err.Error())
		return videos, err
	}
	return videos, nil
}

// QueryVideoByVideoIds query video info by video ids
func QueryVideoByVideoIds(ctx context.Context, videoIds []int64) ([]*Video, error) {
	var videos []*Video
	if err := DB.WithContext(ctx).Where("id in (?)", videoIds).Find(&videos).Error; err != nil {
		klog.Error("QueryVideoByVideoIds error " + err.Error())
		return nil, err
	}
	return videos, nil
}