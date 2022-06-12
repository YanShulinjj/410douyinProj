package controller

import (
	"errors"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

var db *gorm.DB

func init() {
	db, _ = ConnectDataBase()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&VideoSample{})
	db.AutoMigrate(&Comment{})
}

// 连接数据库
func ConnectDataBase() (*gorm.DB, error) {
	dsn := "root:chun@tcp(127.0.0.1:3306)/golang_mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't connect DataBase!")
	}
	return db, nil
}

/***************************** 用户 ********************************/
// 检测用户是否存在
func checkUserName(username string) (User, bool) {
	user := User{}
	result := db.Select("ID").Where("name = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, false
	}
	return user, true
}

// 获得全部用户的简要信息，用于初始化
func GetUsersBriefInfo() []User {
	users := []User{}
	db.Select("ID", "Name", "Password", "LikeVideosID", "FollowCount", "FollowerCount", "FollowID", "FollowerID").Find(&users)
	return users
}

// 从用户信息查找
func FindUserInfo(token string) (User, bool) {
	split := strings.Split(token, "_")
	username, password := split[0], split[1]

	user := User{}
	result := db.Where("name = ? AND password = ? ", username, password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, false
	}
	return user, true
}

// 通过user_id查找用户
func FindUserByID(id uint) (User, bool) {
	user := User{}
	result := db.Select("ID", "Name", "FollowCount", "FollowerCount", "IsFollow", "FollowID", "FollowerID").Where("ID = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, false
	}
	return user, true
}

// 更新用户
func UpdateUser(user User) {
	db.Model(&user).Updates(user)
}

// 添加用户信息
func AddUserInfo(username string, password string) (uint, error) {
	// 向数据库中插入一条数据
	newUser := User{
		Name:     username,
		Password: password,
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		return 0, result.Error
	}
	return newUser.ID, nil
}

// 删除用户
func DeleteUser(token string) error {
	split := strings.Split(token, "_")
	username, password := split[0], split[1]
	user := User{}
	result := db.Where("name = ? AND password = ? ", username, password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	db.Select(clause.Associations).Delete(&user)
	return nil
}

/***************************** 评论 ********************************/
// 添加评论
func AddComment(token string, video_id uint, content string) (Comment, error) {
	// 首先根据token查找到用户ID
	user, exits := FindUserInfo(token)
	if !exits {
		return Comment{}, errors.New("User doesn't exits!")
	}
	timestr := time.Now().Format("2006-01-02 15:04:05")
	new_comment := Comment{
		Content:    content,
		UserId:     user.ID,
		CreateDate: timestr,
	}
	video := VideoSample{}
	db.Where("ID = ?", video_id).First(&video)
	comments := []Comment{}
	result := db.Where("video_refer = ?", video_id).Find(&comments)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Comment{}, result.Error
	}
	comments = append(comments, new_comment)
	db.Model(&video).Update("Comments", comments)

	return new_comment, nil
}

// 删除评论
func DeleteComment(comment_id int) {
	db.Delete(&Comment{}, comment_id)
}

/***************************** 视频 ********************************/
// 添加Video
func AddVideo(token string, videoname string) (VideoSample, error) {
	// 首先根据token查找到用户ID
	user, exits := FindUserInfo(token)
	if !exits {
		return VideoSample{}, errors.New("User doesn't exits!")
	}
	new_video := VideoSample{
		PlayUrl:    BaseURL + "static/" + videoname,
		CoverUrl:   BaseURL + "static/" + videoname + ".jpg",
		IsFavorite: false,
	}
	videos := []VideoSample{}
	result := db.Where("user_refer = ?", user.ID).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return VideoSample{}, result.Error
	}
	videos = append(videos, new_video)
	db.Model(&user).Update("Public", videos)
	// 再次查找
	result = db.Where("user_refer = ?", user.ID).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return VideoSample{}, result.Error
	}

	return videos[len(videos)-1], nil
}

// 删除视频
func DeleteVideo(video_id int) {
	video := FindVideo(video_id)
	db.Select("Comments").Delete(&video)
}

// 查找video
func FindVideo(video_id int) VideoSample {
	video := VideoSample{}
	db.Where("ID = ?", video_id).First(&video)
	return video
}

// 通过user_refer查找video
func FindUsersVideos(user_id uint) []VideoSample {
	videos := []VideoSample{}
	db.Where("user_refer = ?", user_id).Find(&videos)
	return videos
}

// 更新video
func Update(video Video) {
	db.Model(&video).Updates(video)
}

// 返回video列表
func GetVideos() []VideoSample {
	videos := []VideoSample{}
	db.Order("ID desc").Find(&videos)
	return videos
}

// 返回指定视频的comment列表
func GetComments(video_id uint) []Comment {
	comments := []Comment{}
	db.Where("video_refer = ? ", video_id).Find(&comments)
	return comments
}
