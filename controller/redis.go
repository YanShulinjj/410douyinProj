package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

func Connectredis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result() //心跳
	fmt.Println(pong, err)
	return client
}
func Saverediscomment(name string, client *redis.Client, CommentsMap map[uint][]Comment) {

	//info数组转换成json
	inforByte, infoError := json.Marshal(CommentsMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置CommentsMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置CommentsMap成功！")
		}
	}

}
func Getrediscomment(name string, client *redis.Client) map[uint][]Comment {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint][]Comment
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("comment反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint][]Comment{}
}
func SaveVideosBuffer(name string, client *redis.Client, VideosBuffer map[uint]int) {
	inforByte, infoError := json.Marshal(VideosBuffer)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置VideosBuffer出错：", infoErrorStatus)
		} else {
			fmt.Println("设置VideosBuffer成功！")
		}
	}
}
func GetVideosBuffer(name string, client *redis.Client) map[uint]int {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint]int
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("VideosBuffer反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint]int{}
}
func SaveUserFavoriteListMap(name string, client *redis.Client, UserFavoriteListMap map[string][]Video) {
	inforByte, infoError := json.Marshal(UserFavoriteListMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UserFavoriteListMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UserFavoriteListMap成功！")
		}
	}
}
func GetUserFavoriteListMap(name string, client *redis.Client) map[string][]Video {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[string][]Video
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UserFavoriteListMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[string][]Video{}
}
func SaveUsePublishVideosMap(name string, client *redis.Client, UsePublishVideosMap map[uint][]Video) {
	inforByte, infoError := json.Marshal(UsePublishVideosMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UsePublishVideosMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UsePublishVideosMap成功！")
		}
	}
}
func GetUsePublishVideosMap(name string, client *redis.Client) map[uint][]Video {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint][]Video
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UsePublishVideosMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint][]Video{}
}
func SaveUserFollowMap(name string, client *redis.Client, UserFollowMap map[uint][]User) {
	inforByte, infoError := json.Marshal(UserFollowMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UserFollowMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UserFollowMap成功！")
		}
	}
}
func GetUserFollowMap(name string, client *redis.Client) map[uint][]User {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint][]User
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UserFollowMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint][]User{}
}
func SaveUserFollowerMap(name string, client *redis.Client, UserFollowerMap map[uint][]User) {
	inforByte, infoError := json.Marshal(UserFollowerMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UserFollowerMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UserFollowerMap成功！")
		}
	}
}
func GetUserFollowerMap(name string, client *redis.Client) map[uint][]User {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint][]User
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UserFollowerMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint][]User{}
}
func SaveUserFollowCountMap(name string, client *redis.Client, UserFollowCountMap map[uint]int64) {
	inforByte, infoError := json.Marshal(UserFollowCountMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UserFollowCountMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UserFollowCountMap成功！")
		}
	}
}
func GetUserFollowCountMap(name string, client *redis.Client) map[uint]int64 {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint]int64
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UserFollowCountMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint]int64{}
}
func SaveUserFollowerCountMap(name string, client *redis.Client, UserFollowerCountMap map[uint]int64) {
	inforByte, infoError := json.Marshal(UserFollowerCountMap)
	if infoError == nil {
		inforString := string(inforByte)                                        //转换成字符串
		infoErrorStatus := client.Set(name, inforString, 600*time.Second).Err() //设置600秒过期
		if infoErrorStatus != nil {
			fmt.Println("设置UserFollowerCountMap出错：", infoErrorStatus)
		} else {
			fmt.Println("设置UserFollowerCountMap成功！")
		}
	}
}
func GetUserFollowerCountMap(name string, client *redis.Client) map[uint]int64 {
	getInfo, getinfoErr := client.Get(name).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)

	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		var getInfoResult map[uint]int64
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &getInfoResult)

		if unmarsha1Err != nil {
			fmt.Println("UserFollowerCountMap反序列化失败:", unmarsha1Err)
		} else {
			//迭代数据
			for key, value := range getInfoResult {
				fmt.Println(key, value)
			}
			fmt.Println(getInfoResult)
		}
		return getInfoResult
	}
	return map[uint]int64{}
}
