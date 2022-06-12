/* ----------------------------------
*  @author suyame 2022-06-12 15:33:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Addr     string `json:"addr"`
	Port     string `json:"port"`
}

type AddrPort struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}
type MQConfig struct {
	Consumer AddrPort `json:"consumer"`
	Producer AddrPort `json:"producer"`
	Topic    string   `json:"topic"`
	Channel  string   `json:"channel"`
}

type Config struct {
	Mysql   MysqlConfig `json:"mysql_config"`
	BaseURL string      `json:"baseURL"`
	MQ      MQConfig    `json:"message_queue"`
}

var MyConfig Config

func init() {
	readConfig("./config/config.json")
}

func readConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &MyConfig)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return err
	}
	return nil
}
