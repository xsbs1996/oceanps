package oceanps

import (
	"context"
	"github.com/xsbs1996/oceanps/oceanpsfuncs"
	"time"
)

type PublishMsgBody struct {
	Msg     []byte        // 消息体
	Exp     time.Duration // 过期时间
	IsAsync bool          // 是否异步
	Key     string
}

// Publish 发布事件
func (e *EventTopic) Publish(ctx context.Context, body *PublishMsgBody) error {
	err := e.pPManage.PushMsg(ctx, e.queueName, &oceanpsfuncs.PushMsgBody{
		Msg:     body.Msg,
		Exp:     body.Exp,
		IsAsync: body.IsAsync,
		Key:     body.Key,
	})
	if err != nil {
		return err
	}
	return nil
}
