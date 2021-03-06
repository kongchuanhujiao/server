package internal

import (
	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

// setProtocol 设置协议
func setProtocol() {
	client.SystemDeviceInfo.Protocol = client.AndroidPhone
	client.SystemDeviceInfo.Display = []byte("MIRAI.373480.001")
	client.SystemDeviceInfo.Product = []byte("mirai")
	client.SystemDeviceInfo.Device = []byte("mirai")
	client.SystemDeviceInfo.Board = []byte("mirai")
	client.SystemDeviceInfo.Model = []byte("mirai")
	client.SystemDeviceInfo.FingerPrint = []byte("mamoe/mirai/mirai:10/MIRAI.200122.001/6671789:user/release-keys")
	client.SystemDeviceInfo.BootId = []byte("7794a02c-d854-18ac-649e-35fedfd0b37a")
	client.SystemDeviceInfo.ProcVersion = []byte("Linux version 3.0.31-47Fxpwhn (android-build@xxx.xxx.xxx.xxx.com)")
	client.SystemDeviceInfo.IMEI = "678319144775066"
}

// setLogger 设置日志打印
func setLogger(_ *client.QQClient, e *client.LogEvent) {
	switch e.Type {
	case "INFO":
		loggerr.Info("协议信息：" + e.Message)
	case "ERROR":
		loggerr.Error("协议错误：" + e.Message)
	}
}

// login 登录
func (q *QQ) login() (err error) {

	for res, err := q.client.Login(); err != nil || !res.Success; res, err = q.client.Login() {

		if err != nil {
			if err == client.ErrAlreadyOnline {
				break
			}

			loggerr.Error("登录失败", zap.Error(err))
			return err
		}

		loggerr.Panic("无法登录：" + res.ErrorMessage)
	}

	err = q.client.ReloadGroupList()
	if err != nil {
		loggerr.Error("加载群列表失败", zap.Error(err))
		return err
	}

	err = q.client.ReloadFriendList()
	if err != nil {
		loggerr.Error("加载好友列表失败", zap.Error(err))
		return err
	}

	loggerr.Info("登录成功：" + q.client.Nickname)

	return
}

// setEventHandle 设置事件处理器
func (q *QQ) setEventHandle() {

	q.client.OnGroupMessage(q.receiveGroupMessage)
	q.client.OnPrivateMessage(q.receiveFriendMessage)

	// 更新服务器
	q.client.OnServerUpdated(func(_ *client.QQClient, e *client.ServerUpdatedEvent) bool {
		loggerr.Warn("更新服务器", zap.Any("信息", e.Servers))
		return true
	})

	// 断线重连
	q.client.OnDisconnected(func(_ *client.QQClient, e *client.ClientDisconnectedEvent) {
		for {

			loggerr.Warn("连接丢失，重连中...")
			if err := q.login(); err != nil {
				loggerr.Warn("重登录失败，再次尝试中...", zap.Error(err))
				continue
			}

			return
		}
	})
}
