package qq

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

var loggerr = logger.Named("QQ客户端") // loggerr 日志

// QQ QQ 客户端
type QQ struct {
	client   *client.QQClient       // 客户端
	callback clientspublic.Callback // 回调
}

// NewQQClient 新建 QQ 客户端
func NewQQClient(a uint64, p string) (q *QQ) {

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

func (q *QQ) SendMessage(m *clientspublic.Message) {
	ms := q.transformToMiraiGO(m)
	if m.Target.Group == nil {
		q.client.SendGroupMessage(int64(m.Target.Group.ID), ms)
	} else {
		q.client.SendPrivateMessage(int64(m.Target.ID), ms)
	}
	loggerr.Info("发送消息", zap.Any("消息", m.Chain))
}

func (q *QQ) ReceiveMessage(m *clientspublic.Message) {
	logger.Debug("接收消息", zap.Any("消息", m))
	if q.callback != nil {
		q.callback(m)
	}
}

func (q *QQ) SetCallback(f clientspublic.Callback) {
	q.callback = f
}
