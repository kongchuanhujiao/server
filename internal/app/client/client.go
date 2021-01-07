package client

import (
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/client/internal"
	"coding.net/kongchuanhujiao/server/internal/pkg/configs"
)

var cli *internal.QQ // client 客户端

// NewClients 新建客户端
func NewClients() {
	conf := configs.GetConfigs()
	cli = internal.NewClient(conf.QQNumber, conf.QQPassword)
}

// GetClient 获取客户端。
// 执行函数：NewClients 前调用返回值为 nil
func GetClient() *internal.QQ { return cli }

// SetCallback 设置回调
func SetCallback(f clientmsg.Callback) { cli.SetCallback(f) }