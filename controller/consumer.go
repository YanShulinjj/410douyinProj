/* ----------------------------------
*  @author suyame 2022-06-11 16:43:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package controller

import (
	"encoding/json"
	"github.com/RaymondCode/simple-demo/mylog"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

func RunComsumer() {
	topic := MyConfig.MQ.Topic
	channel := MyConfig.MQ.Channel
	addr := MyConfig.MQ.Consumer.Addr + ":" + MyConfig.MQ.Consumer.Port
	// err := initConsumer("test1", "test-channel1", "127.0.0.1:4161")
	err := initConsumer(topic, channel, addr)
	if err != nil {
		log.Fatal("init Consumer error")
	}
	select {}
}

type nsqHandler struct {
	nsqConsumer      *nsq.Consumer
	messagesReceived int
}

func initConsumer(topic, channel, addr string) error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 3 * time.Second
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Println("init Consumer NewConsumer error:", err)
		return err
	}

	handler := &nsqHandler{nsqConsumer: c}
	c.AddHandler(handler)

	err = c.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Println("init Consumer ConnectToNSQLookupd error:", err)
		return err
	}
	return nil
}

// 处理消息
func (nh *nsqHandler) HandleMessage(msg *nsq.Message) error {
	nh.messagesReceived++
	// fmt.Printf("receive ID:%s,addr:%s,message:%s", msg.ID, msg.NSQDAddress, string(msg.Body))
	// fmt.Println()
	qms := MQmessage{}
	err := json.Unmarshal(msg.Body, &qms)
	if err != nil {
		mylog.Logger.Fatal("mqmessage unmarshall error!", err)
		return err
	}
	Interaction(qms)
	return nil
}
