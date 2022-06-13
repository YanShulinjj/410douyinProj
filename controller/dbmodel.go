package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/mylog"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
	"time"
)

var db *gorm.DB
var Client *redis.Client

func init() {
	db, _ = ConnectDataBase()
	Client = ConnectRedis()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&VideoSample{})
	db.AutoMigrate(&Comment{})
}

// 连接数据库
func ConnectDataBase() (*gorm.DB, error) {

	user := MyConfig.Mysql.User
	password := MyConfig.Mysql.Password
	addr := MyConfig.Mysql.Addr
	port := MyConfig.Mysql.Port
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/golang_mysql?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, addr, port)
	// dsn := "root:19990221@tcp(127.0.0.1:3306)/golang_mysql?charset=utf8mb4&parseTime=True&loc=Local"
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



//判断账号是否为邮箱格式

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
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
	// 通过发送消息来更新
	// db.Model(&user).Updates(user)
	Public(MQmessage{
		DataType: 0,
		OpType:   1,
		Data:     user,
	})
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
	video.CommentCount++
	video.Comments = &comments
	// db.Model(&video).Updates(video)
	// db.Model(&video).Update("Comments", comments)
	// 发送更新的消息
	Public(MQmessage{
		DataType: 2,
		OpType:   1,
		Data:     video,
	})

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
		PlayUrl:    MyConfig.BaseURL + "static/" + videoname,
		CoverUrl:   MyConfig.BaseURL + "static/" + videoname + ".jpg",
		IsFavorite: false,
	}
	videos := []VideoSample{}
	result := db.Where("user_refer = ?", user.ID).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return VideoSample{}, result.Error
	}
	videos = append(videos, new_video)
	db.Model(&user).Update("Public", videos)
	// 通过发送消息延迟更新
	// user.Public = &videos
	// Public(MQmessage{
	// 	DataType: 0,
	// 	OpType:   1,
	// 	Data:     user,
	// })
	// 再次查找
	result = db.Where("user_refer = ?", user.ID).Find(&videos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return VideoSample{}, result.Error
	}
	fmt.Println("Adding videos... ", videos)
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
	// 通过发送消息更新
	Public(MQmessage{
		DataType: 2,
		OpType:   1,
		Data:     video,
	})
	// db.Model(&video).Updates(video)
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

/*** 解析消息队列中的消息并执行对应操作 ****/
/**
DataType: 0 User, 1 Comment, 2 Video
OPType:   0 Add,  1  Modify, 2 Delete
*/

func crudUser(user User, optype uint) error {
	switch optype {
	case 0: // Add
		result := db.Create(&user)
		if result.Error != nil {
			mylog.Logger.Printf("Add user [user_id = %d] error, %s\n", user.ID, result.Error)
			return result.Error
		}
	case 1: // modify
		db.Model(&user).Updates(user)
	case 2: // delete
		db.Select(clause.Associations).Delete(&user)
	default:
		mylog.Logger.Println("Unknown Optype!")
	}
	return nil
}

func crudComment(comment Comment, optype uint) error {
	switch optype {
	case 0: // Add
		result := db.Create(&comment)
		if result.Error != nil {
			mylog.Logger.Printf("Add comment [comment_id = %d] error, %s\n", comment.ID, result.Error)
			return result.Error
		}
	case 1: // modify
		db.Model(&comment).Updates(comment)
	case 2: // delete
		db.Select(clause.Associations).Delete(&comment)
	default:
		mylog.Logger.Println("Unknown Optype!")
	}
	return nil
}

func crudVideo(video VideoSample, optype uint) error {
	switch optype {
	case 0: // Add
		result := db.Create(&video)
		if result.Error != nil {
			mylog.Logger.Printf("Add video [video_id = %d] error, %s\n", video.ID, result.Error)
			return result.Error
		}
	case 1: // modify
		db.Model(&video).Updates(video)
	case 2: // delete
		db.Select(clause.Associations).Delete(&video)
	default:
		mylog.Logger.Println("Unknown Optype!")
	}
	return nil
}

func Interaction(qmessage MQmessage) {
	// 根据队列中的消息区分出对什么表进行什么操作
	switch qmessage.DataType {
	case 0: // User
		user := User{}
		arr, err := json.Marshal(qmessage.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(arr, &user)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = crudUser(user, qmessage.OpType)
		if err != nil {
			mylog.Logger.Panicf("user[user_id=%d], Op[optype=%d], error!, %s", user.ID, qmessage.OpType, err)
		}
	case 1: // Comment
		comment := Comment{}
		arr, err := json.Marshal(qmessage.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(arr, &comment)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = crudComment(comment, qmessage.OpType)
		if err != nil {
			mylog.Logger.Panicf("comment[comment_id=%d], Op[optype=%d], error!, %s", comment.ID, qmessage.OpType, err)
		}
	case 2: // Video
		video := VideoSample{}
		arr, err := json.Marshal(qmessage.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(arr, &video)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = crudVideo(video, qmessage.OpType)
		if err != nil {
			mylog.Logger.Panicf("video[video_id=%d], Op[optype=%d], error!, %s", video.ID, qmessage.OpType, err)
		}
	default:
	}
}
