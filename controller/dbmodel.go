package controller

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

var db *gorm.DB

func init() {
	db, _ = ConnectDataBase()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&VideoSample{})
	db.AutoMigrate(&Comment{})
}

// Init DataSet
func ConnectDataBase() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:1002@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                             // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                            // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                            // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                            // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                           // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true, //建表的时候不会建立物理外键。 主张逻辑外键 (代码里面自动体现外键关系)
	})
	if err != nil {
		fmt.Println(err)
		panic("Can't connect DataBase!")
	}
	return db, nil
}

// 查找用户名

func FindUserName(username string) (User, bool) {
	user := User{}
	result := db.Where("name = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, false
	}
	// 对用户名的密码掩码
	user.Password = "******"
	return user, true
}

// 从用户信息查找

func FindUserInfo(token string) (User, bool) {

	fmt.Println("TOKEN: ", token)
	split := strings.Split(token, "_")
	fmt.Println("Split: ", split)
	username, password := split[0], split[1]
	fmt.Println(username, password)

	user := User{}
	result := db.Where("name = ? AND password = ? ", username, password).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, false
	}
	return user, true
}

// 通过id查找用户

func FindUserByID(id uint) (User, bool) {
	user := User{}
	result := db.Select("ID", "Name", "FollowCount", "FollowerCount", "IsFollow").Where("ID = ?", id).First(&user)
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
	// db.AutoMigrate(&User{})
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
	split := strings.Split(token, "&")
	username, password := split[0], split[1]
	user := User{}
	result := db.Where("name = ? AND password = ? ", username, password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	db.Select(clause.Associations).Delete(&user)
	return nil
}

// 添加评论

func AddComment(token string, video_id uint, content string) (Comment, error) {
	db.AutoMigrate(&Comment{})
	// 首先根据token查找到用户ID
	user, exits := FindUserInfo(token)
	if !exits {
		return Comment{}, errors.New("User doesn't exits!")
	}
	new_comment := Comment{
		Content: content,
		UserId:  user.ID,
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

//
func GetUsersBriefInfo() []User {
	users := []User{}
	db.Select("ID", "Name", "Password", "LikeVideosID", "FollowCount", "FollowerCount", "FollowID", "FollowerID").Find(&users)
	return users
}

// func main() {
// 	// create
// 	username := "junjun"
// 	password := "19980419"
// 	AddUserInfo(username, password)

// 	// 查找user
// 	user, exits := FindUserInfo(username, password)
// 	if exits {
// 		fmt.Println(user)
// 	}
// 	token := username + "&" + password
// 	// 发布视频
// 	AddVideo(token, "bear.mp4")
// 	// 发布评论
// 	AddComment(token, int64(1), "这个视频好好看！")
// 	AddComment(token, int64(1), "this video pretty good!")
// 	// // 删除评论
// 	// DeleteVideo(1)
// 	// 删除用户
// 	// err := DeleteUser(token)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	videos := GetVideos()
// 	for _, video := range videos {
// 		fmt.Println(video.UserId)
// 		fmt.Println(video.UserId)
// 	}
// }
