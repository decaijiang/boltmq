package registry

import (
	"git.oschina.net/cloudzone/smartgo/stgbroker/client"
	"git.oschina.net/cloudzone/smartgo/stgnet/netm"
)

// BrokerHousekeepingServices Broker活动检测服务
//
// (1)ChannelEventListener是RocketMQ封装Netty向外暴露的一个接口层
// (2)NameSrv监测Broker的死亡：当Broker和NameSrv之间的长连接断掉之后，后续的ChannelEventListener里面的函数就会被回调，从而触发NameServer的路由信息更新
//
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
type BrokerHousekeepingService struct {
	NamesrvController *DefaultNamesrvController
}

// NewBrokerHousekeepingService 初始化Broker活动检测服务
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
func NewBrokerHousekeepingService(controller *DefaultNamesrvController) client.ChannelEventListener {
	brokerHousekeepingService := &BrokerHousekeepingService{
		NamesrvController: controller,
	}
	return brokerHousekeepingService
}

// onChannelConnect
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
func (self *BrokerHousekeepingService) OnChannelConnect(ctx netm.Context) {

}

// onChannelClose Channel被关闭,通知Topic路由管理器，清除无效Broker
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
func (self *BrokerHousekeepingService) OnChannelClose(ctx netm.Context) {
	self.NamesrvController.RouteInfoManager.onChannelDestroy(ctx.RemoteAddr().String(), ctx)
}

// onChannelException Channel出现异常,通知Topic路由管理器，清除无效Broker
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
func (self *BrokerHousekeepingService) OnChannelException(ctx netm.Context) {
	self.NamesrvController.RouteInfoManager.onChannelDestroy(ctx.RemoteAddr().String(), ctx)
}

// onChannelIdle Channe的Idle时间超时,通知Topic路由管理器，清除无效Brokers
// Author: tianyuliang, <tianyuliang@gome.com.cn>
// Since: 2017/9/6
func (self *BrokerHousekeepingService) OnChannelIdle(ctx netm.Context) {
	self.NamesrvController.RouteInfoManager.onChannelDestroy(ctx.RemoteAddr().String(), ctx)
}
