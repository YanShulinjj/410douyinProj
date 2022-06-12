/* ----------------------------------
*  @author suyame 2022-06-11 16:26:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
)

type MQmessage struct {
	DataType uint        `json:"data_type,omitempty"`
	OpType   uint        `json:"op_type,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type nsqProducer struct {
	*nsq.Producer
}

var producer *nsqProducer

func init() {
	// strIP1 := "127.0.0.1:4150"
	strIP1 := MyConfig.MQ.Producer.Addr + ":" + MyConfig.MQ.Producer.Port
	producer, _ = initProducer(strIP1)
}

func Public(message MQmessage) error {
	mess_bytes, err := json.Marshal(message)
	// fmt.Println(string(mess_bytes))
	if err != nil {
		fmt.Println(err)
	}
	topic := MyConfig.MQ.Topic
	err = producer.Publish(topic, mess_bytes)
	if err != nil {
		log.Fatal("producer1 public error:", err)
		return err
	}
	return nil
}

// 初始化生产者
func initProducer(addr string) (*nsqProducer, error) {
	// fmt.Println("init producer address:", addr)
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &nsqProducer{producer}, nil
}
