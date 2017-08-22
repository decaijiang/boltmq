package process

import (
	"git.oschina.net/cloudzone/smartgo/stgcommon/protocol/heartbeat"
	"git.oschina.net/cloudzone/smartgo/stgcommon/protocol/route"
	"git.oschina.net/cloudzone/smartgo/stgcommon/message"
	"git.oschina.net/cloudzone/smartgo/stgcommon/protocol/header"
	"strings"
	"git.oschina.net/cloudzone/smartgo/stgclient"
	"git.oschina.net/cloudzone/smartgo/stgcommon/logger"
	"git.oschina.net/cloudzone/smartgo/stgclient/consumer"
	"git.oschina.net/cloudzone/smartgo/stgnet/protocol"
	cprotocol "git.oschina.net/cloudzone/smartgo/stgcommon/protocol"
	"git.oschina.net/cloudzone/smartgo/stgcommon"
	"errors"
)

// MQClientAPIImpl: 内部使用核心处理api
// Author: yintongqiang
// Since:  2017/8/8

type MQClientAPIImpl struct {
	ClientRemotingProcessor *ClientRemotingProcessor
	ProjectGroupPrefix      string
}

func NewMQClientAPIImpl(clientRemotingProcessor *ClientRemotingProcessor) *MQClientAPIImpl {

	return &MQClientAPIImpl{
		ClientRemotingProcessor:clientRemotingProcessor,
	}
}
// 调用romoting的start
func (impl *MQClientAPIImpl)Start() {
	//todo 初始化远程
	//impl.ProjectGroupPrefix

}
// 关闭romoting
func (impl *MQClientAPIImpl)Shutdwon() {
	//todo 关闭远程

}
// 发送心跳到broker
func (impl *MQClientAPIImpl)sendHeartbeat(addr string, heartbeatData *heartbeat.HeartbeatData, timeoutMillis int64) {
	if strings.EqualFold(impl.ProjectGroupPrefix, "") {
		consumerDatas := heartbeatData.ConsumerDataSet
		for data := range consumerDatas.Iterator().C {
			consumerData := data.(heartbeat.ConsumerData)
			consumerData.GroupName = stgclient.BuildWithProjectGroup(consumerData.GroupName, impl.ProjectGroupPrefix)
			subscriptionDatas := consumerData.SubscriptionDataSet
			for subData := range subscriptionDatas.Iterator().C {
				subscriptionData := subData.(heartbeat.SubscriptionData)
				subscriptionData.Topic = stgclient.BuildWithProjectGroup(subscriptionData.Topic, impl.ProjectGroupPrefix)
			}
		}
		producerDatas := heartbeatData.ProducerDataSet
		for pData := range producerDatas.Iterator().C {
			producerData := pData.(*heartbeat.ProducerData)
			producerData.GroupName = stgclient.BuildWithProjectGroup(producerData.GroupName, impl.ProjectGroupPrefix)
		}
	}
	//todo 创建request调用invokeSync
}

func (impl *MQClientAPIImpl)GetDefaultTopicRouteInfoFromNameServer(topic string, timeoutMillis int64) *route.TopicRouteData {
	topicWithProjectGroup := topic
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		topicWithProjectGroup = stgclient.BuildWithProjectGroup(topic, impl.ProjectGroupPrefix)
	}
	requestHeader:=header.GetRouteInfoRequestHeader{Topic:topicWithProjectGroup}
	request:=protocol.CreateRequestCommand(cprotocol.GET_ROUTEINTO_BY_TOPIC,&requestHeader)
	logger.Infof(request.RemotingVersionKey)
	//todo 调用远程生成
	reponse:=protocol.CreateResponseCommand()
	switch reponse.Code {
	case cprotocol.TOPIC_NOT_EXIST:
		logger.Warnf("get Topic [%v] RouteInfoFromNameServer is not exist value", topic)
	case cprotocol.SUCCESS:
		body:=reponse.Body
		if len(body)>0{
			//todo
		}
	}
	return &route.TopicRouteData{}
}

func (impl *MQClientAPIImpl)GetTopicRouteInfoFromNameServer(topic string, timeoutMillis int64) *route.TopicRouteData {

	topicWithProjectGroup:=topic
	if !strings.EqualFold(impl.ProjectGroupPrefix,""){
		topicWithProjectGroup=stgclient.BuildWithProjectGroup(topic,impl.ProjectGroupPrefix)
	}
	requestHeader:=&header.GetRouteInfoRequestHeader{Topic:topicWithProjectGroup}
	request:=protocol.CreateRequestCommand(cprotocol.GET_ROUTEINTO_BY_TOPIC,requestHeader)
	logger.Info(request.RemotingVersionKey)
    //todo response处理
	routeData := &route.TopicRouteData{}
	routeData.QueueDatas = append(routeData.QueueDatas, &route.QueueData{BrokerName:"broker-master2", ReadQueueNums:8, WriteQueueNums:8, Perm:6, TopicSynFlag:0})
	mapBrokerAddrs := make(map[int]string)
	mapBrokerAddrs[0] = "10.128.31.124:10911"
	mapBrokerAddrs[1] = "10.128.31.125:10911"
	routeData.BrokerDatas = append(routeData.BrokerDatas, &route.BrokerData{BrokerName:"broker-master2", BrokerAddrs:mapBrokerAddrs})
	return routeData
}

func (impl *MQClientAPIImpl)SendMessage(addr string, brokerName string, msg *message.Message, requestHeader header.SendMessageRequestHeader,
timeoutMillis int64, communicationMode CommunicationMode, sendCallback SendCallback) (SendResult,error) {
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		msg.Topic = stgclient.BuildWithProjectGroup(msg.Topic, impl.ProjectGroupPrefix)
		requestHeader.ProducerGroup = stgclient.BuildWithProjectGroup(requestHeader.ProducerGroup, impl.ProjectGroupPrefix)
		requestHeader.Topic = stgclient.BuildWithProjectGroup(requestHeader.Topic, impl.ProjectGroupPrefix)
	}
	// 默认send采用v2版本
	requestHeaderV2 := header.CreateSendMessageRequestHeaderV2(&requestHeader)
	request := protocol.CreateRequestCommand(cprotocol.SEND_MESSAGE_V2, requestHeaderV2)
	switch (communicationMode) {
	case ONEWAY:
	case ASYNC:
	case SYNC:
		return impl.sendMessageSync(addr, brokerName, msg, timeoutMillis, request)
	default:
		break
	}

	return SendResult{},errors.New("SendMessage error")
}

func (impl *MQClientAPIImpl)sendMessageSync(addr string, brokerName string, msg *message.Message, timeoutMillis int64, request *protocol.RemotingCommand) (SendResult,error) {
	//todo 调用远程生成response
	response := protocol.CreateResponseCommand()
	return impl.processSendResponse(brokerName, msg, response)
}

// 处理发送消息响应
func (impl *MQClientAPIImpl)processSendResponse(brokerName string, msg *message.Message, response *protocol.RemotingCommand) (SendResult,error) {
	switch response.Code {
	case cprotocol.FLUSH_DISK_TIMEOUT:
	case cprotocol.FLUSH_SLAVE_TIMEOUT:
	case cprotocol.SLAVE_NOT_AVAILABLE:
		logger.Warnf("brokerName %v SLAVE_NOT_AVAILABLE", brokerName)
	case cprotocol.SUCCESS:
		sendStatus := SEND_OK
		switch response.Code {
		case cprotocol.FLUSH_DISK_TIMEOUT:
			sendStatus = FLUSH_DISK_TIMEOUT
		case cprotocol.FLUSH_SLAVE_TIMEOUT:
			sendStatus = FLUSH_SLAVE_TIMEOUT
		case cprotocol.SLAVE_NOT_AVAILABLE:
			sendStatus = SLAVE_NOT_AVAILABLE
		case cprotocol.SUCCESS:
			sendStatus = SEND_OK
		default:
		}
		//todo 需从responde中解析出responseHeader
		responseHeader := header.SendMessageResponseHeader{}
		messageQueue := message.MessageQueue{Topic:msg.Topic, BrokerName:brokerName, QueueId:responseHeader.QueueId}
		sendResult := NewSendResult(sendStatus, responseHeader.MsgId, messageQueue, responseHeader.QueueOffset, impl.ProjectGroupPrefix)
		sendResult.TransactionId = responseHeader.TransactionId
		return sendResult,nil
	}
	return SendResult{},errors.New("processSendResponse error")
}

func (impl *MQClientAPIImpl)UpdateConsumerOffsetOneway(addr string, requestHeader header.UpdateConsumerOffsetRequestHeader, timeoutMillis int64) {
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		requestHeader.ConsumerGroup = stgclient.BuildWithProjectGroup(requestHeader.ConsumerGroup, impl.ProjectGroupPrefix)
		requestHeader.Topic = stgclient.BuildWithProjectGroup(requestHeader.Topic, impl.ProjectGroupPrefix)
	}
	// todo 创建request
}

func (impl *MQClientAPIImpl)GetConsumerIdListByGroup(addr string, consumerGroup string, timeoutMillis int64) []string {
	consumerGroupWithProjectGroup := consumerGroup
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		consumerGroupWithProjectGroup = stgclient.BuildWithProjectGroup(consumerGroup, impl.ProjectGroupPrefix)
	}
	requestHeader := header.GetConsumerListByGroupRequestHeader{ConsumerGroup:consumerGroupWithProjectGroup}
	logger.Infof(requestHeader.ConsumerGroup)
	// todo 创建request
	return []string{}
}

func (impl *MQClientAPIImpl)GetMaxOffset(addr string, topic string, queueId int, timeoutMillis int64) int64 {
	topicWithProjectGroup := topic
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		topicWithProjectGroup = stgclient.BuildWithProjectGroup(topic, impl.ProjectGroupPrefix)
	}
	logger.Infof(topicWithProjectGroup)
	// todo 创建request
	return 10
}

func (impl *MQClientAPIImpl)PullMessage(addr string, requestHeader header.PullMessageRequestHeader,
timeoutMillis int, communicationMode CommunicationMode, pullCallback consumer.PullCallback) consumer.PullResult {
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		requestHeader.ConsumerGroup = stgclient.BuildWithProjectGroup(requestHeader.ConsumerGroup, impl.ProjectGroupPrefix)
		requestHeader.Topic = stgclient.BuildWithProjectGroup(requestHeader.Topic, impl.ProjectGroupPrefix)
	}
	request := protocol.CreateRequestCommand(cprotocol.PULL_MESSAGE, &requestHeader)
	switch communicationMode {
	case ONEWAY:
	case ASYNC:
		impl.pullMessageAsync(addr, request, timeoutMillis, pullCallback)
	case SYNC:
	default:

	}
	return consumer.PullResult{}
}

func (impl *MQClientAPIImpl)UpdateNameServerAddressList(addrs string) {
	strings.Split(addrs, ";")
	//todo 设置remotingclient的nameserver
}

func (impl *MQClientAPIImpl)consumerSendMessageBack(addr string,msg message.MessageExt,consumerGroup string,delayLevel int,timeoutMillis int) {
	consumerGroupWithProjectGroup := consumerGroup
	if !strings.EqualFold(impl.ProjectGroupPrefix, "") {
		consumerGroupWithProjectGroup = stgclient.BuildWithProjectGroup(consumerGroup, impl.ProjectGroupPrefix)
		msg.Topic=stgclient.BuildWithProjectGroup(msg.Topic,impl.ProjectGroupPrefix)
	}
	requestHeader:=header.ConsumerSendMsgBackRequestHeader{
    Group:consumerGroupWithProjectGroup,
		OriginTopic:msg.Topic,
		Offset:msg.CommitLogOffset,
		DelayLevel:delayLevel,
		OriginMsgId:msg.MsgId,
	}
	request:=protocol.CreateRequestCommand(cprotocol.CONSUMER_SEND_MSG_BACK,&requestHeader)
	logger.Infof(request.RemotingVersionKey)
	//todo remotingclient invokeSync
}


func (impl *MQClientAPIImpl)pullMessageAsync(addr string, request *protocol.RemotingCommand, timeoutMillis int, pullCallback consumer.PullCallback) {
  //todo
}

func (impl *MQClientAPIImpl)unRegisterClient(addr,clientID,producerGroup,consumerGroup string, timeoutMillis int) {
	 producerGroupWithProjectGroup:= producerGroup
	 consumerGroupWithProjectGroup:= consumerGroup
	if !strings.EqualFold("",impl.ProjectGroupPrefix){
		producerGroupWithProjectGroup=stgclient.BuildWithProjectGroup(producerGroup,impl.ProjectGroupPrefix)
		consumerGroupWithProjectGroup=stgclient.BuildWithProjectGroup(consumerGroup,impl.ProjectGroupPrefix)
	}
	requestHeader:=header.UnregisterClientRequestHeader{ClientID:clientID,ProducerGroup:producerGroupWithProjectGroup,ConsumerGroup:consumerGroupWithProjectGroup}
    request:=protocol.CreateRequestCommand(cprotocol.UNREGISTER_CLIENT,&requestHeader)
	logger.Infof(request.RemotingVersionKey)
	//todo invokeSync
}


func (impl *MQClientAPIImpl)processPullResponse(response protocol.RemotingCommand) *PullResultExt {
	pullStatus := consumer.NO_NEW_MSG
	switch response.Code {
	case cprotocol.SUCCESS:
		pullStatus = consumer.FOUND
	case cprotocol.PULL_NOT_FOUND:
		pullStatus = consumer.NO_NEW_MSG
	case cprotocol.PULL_RETRY_IMMEDIATELY:
		pullStatus = consumer.NO_MATCHED_MSG
	case cprotocol.PULL_OFFSET_MOVED:
		pullStatus = consumer.OFFSET_ILLEGAL
	}
	//todo response生产PullMessageResponseHeader
	reponseHeader := &header.PullMessageResponseHeader{}
	return NewPullResultExt(pullStatus, reponseHeader.NextBeginOffset,
		reponseHeader.MaxOffset, reponseHeader.MaxOffset, nil, reponseHeader.SuggestWhichBrokerId, response.Body)
}
// 创建topic
func (impl *MQClientAPIImpl)CreateTopic(addr,defaultTopic string,topicConfig stgcommon.TopicConfig,timeoutMillis int){
	topicWithProjectGroup:=topicConfig.TopicName
	if !strings.EqualFold("",impl.ProjectGroupPrefix){
		topicWithProjectGroup=stgclient.BuildWithProjectGroup(topicConfig.TopicName,impl.ProjectGroupPrefix)
	}
	requestHeader:=header.CreateTopicRequestHeader{Topic:topicWithProjectGroup,
	DefaultTopic:defaultTopic,ReadQueueNums:topicConfig.ReadQueueNums,WriteQueueNums:topicConfig.WriteQueueNums,
		TopicFilterType:topicConfig.TopicFilterType.String(),TopicSysFlag:topicConfig.TopicSysFlag,Order:topicConfig.Order,
	    Perm:topicConfig.Perm}
	request:=protocol.CreateRequestCommand(cprotocol.UPDATE_AND_CREATE_TOPIC,&requestHeader)
	logger.Infof(request.RemotingVersionKey)
	//todo remoting invoke
}