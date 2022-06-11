/* ----------------------------------
*  @author suyame 2022-06-10 13:25:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package mylog

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	logFile, err := os.OpenFile("./mylog/a.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	Logger = log.New(logFile, "[410proj]", log.Ldate|log.Ltime|log.Lshortfile)
}
