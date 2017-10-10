package main

import (
	"fmt"
	"git.oschina.net/cloudzone/smartgo/stgclient/consumer"
	"git.oschina.net/cloudzone/smartgo/stgclient/consumer/listener"
	"git.oschina.net/cloudzone/smartgo/stgclient/process"
	"git.oschina.net/cloudzone/smartgo/stgcommon/message"
	"git.oschina.net/cloudzone/smartgo/stgcommon/protocol/heartbeat"
	"time"
)

type MessageListenerImpl struct {
}

func (listenerImpl *MessageListenerImpl) ConsumeMessage(msgs []*message.MessageExt, context *consumer.ConsumeConcurrentlyContext) listener.ConsumeConcurrentlyStatus {
	for _, msg := range msgs {
		//fmt.Println(msg)
		fmt.Println(string(msg.Body))
	}
	return listener.CONSUME_SUCCESS
}
func taskC() {
	t := time.NewTicker(time.Second * 1000)
	for {
		select {
		case <-t.C:
		}

	}
}

func main() {
	defaultMQPushConsumer := process.NewDefaultMQPushConsumer("myConsumerGroup")
	defaultMQPushConsumer.SetConsumeFromWhere(heartbeat.CONSUME_FROM_LAST_OFFSET)
	defaultMQPushConsumer.SetMessageModel(heartbeat.CLUSTERING)
	defaultMQPushConsumer.SetNamesrvAddr("10.112.68.189:9876")
	defaultMQPushConsumer.Subscribe("cloudzone10", "tagA")
	defaultMQPushConsumer.RegisterMessageListener(&MessageListenerImpl{})
	defaultMQPushConsumer.Start()
	time.Sleep(time.Second * 6000)
	defaultMQPushConsumer.Shutdown()
	go taskC()
	select {}
}