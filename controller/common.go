package controller

import (
	"gorm.io/gorm"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type Video struct {
	Id            uint   `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type VideoSample struct {
	ID            uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	UpdatedAt     time.Time      `json:"created_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PlayUrl       string         `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string         `json:"cover_url,omitempty"`
	FavoriteCount int64          `json:"favorite_count,omitempty"`
	CommentCount  int64          `json:"comment_count,omitempty"`
	Comments      *[]Comment     `gorm:"foreignKey:VideoRefer"`
	IsFavorite    bool           `json:"is_favorite,omitempty"`
	UserRefer     uint           `json:"author"`
}

type Comment struct {
	ID         uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"created_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Content    string         `json:"content,omitempty"`
	CreateDate string         `json:"create_date,omitempty"`
	VideoRefer uint           `json:"video_refer,omitempty"`
	UserId     uint           `json:"user_id,omitempty"`
}

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	UpdatedAt     time.Time      `json:"created_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name          string         `json:"name,omitempty"`
	Password      string         `json:"password,omitempty"`
	FollowCount   int64          `json:"follow_count,omitempty"`
	FollowerCount int64          `json:"follower_count,omitempty"`
	FollowID      string         `json:"follow_id,omitempty"`
	FollowerID    string         `json:"follower_id,omitempty"`
	IsFollow      bool           `json:"is_follow,omitempty"`
	LikeVideosID  string         `json:"like_videos,omitempty"`
	Public        *[]VideoSample `gorm:"foreignKey:UserRefer" json:"video,omitempty"`
}
