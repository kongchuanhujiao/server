package internal

import (
	"coding.net/kongchuanhujiao/server/internal/app/client/message"
	wenda2 "coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

var loggerr = logger.Named("QQ") // loggerr 日志

// QQ 客户端
type QQ struct {
	client   *client.QQClient // 客户端
	callback message.Callback // 回调
}

// NewClient 新建 QQ 客户端
func NewClient(a uint64, p string) (q *QQ) {

	setProtocol()
	c := client.NewClient(int64(a), p)
	c.OnLog(setLogger)

	// 读取配置信息
	q = &QQ{client: c}

	if err := q.login(); err != nil {
		loggerr.Panic("登录失败", zap.Error(err))
	}
	q.setEventHandle()

	return
}

// SendMessage 发送消息
func (q *QQ) SendMessage(m *message.Message) {
	ms := q.transformToMiraiGO(m)
	if m.Target.Group != nil {
		q.client.SendGroupMessage(int64(m.Target.Group.ID), ms)
	} else {
		q.client.SendPrivateMessage(int64(m.Target.ID), ms)
	}
}

// ReceiveMessage 接收消息
func (q *QQ) ReceiveMessage(m *message.Message) {
	loggerr.Debug("接收消息", zap.Any("消息", m))
	if q.callback != nil && len(m.Chain) != 0 {
		q.callback(m)
	}
}

// SetCallback 设置回调
func (q *QQ) SetCallback(f message.Callback) {
	q.callback = f
}

// GetGroups 获取群
func (q *QQ) GetGroups() *wenda2.Groups {
	g := wenda2.Groups{}
	for _, v := range q.client.GroupList {
		g[uint64(v.Code)] = v.Name
	}
	return &g
}

// GetGroupName 获取群名称
func (q *QQ) GetGroupName(i uint64) string { return q.client.FindGroup(int64(i)).Name }

// GetGroupMembers 获取群成员
func (q *QQ) GetGroupMembers(i uint64) *wenda2.GroupMembers {
	m := wenda2.GroupMembers{}
	for _, v := range q.client.FindGroup(int64(i)).Members {
		m[uint64(v.Uin)] = v.DisplayName()
	}
	return &m
}
