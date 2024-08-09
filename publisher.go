package oceanps

import (
	"context"
	"time"
)

// Publish 发布事件
func (e *EventTopic) Publish(ctx context.Context, msg []byte, exp time.Duration, extra ...string) error {
	err := e.pPManage.PushMsgFn(ctx, e.queueName, msg, exp, extra...)
	if err != nil {
		return err
	}
	return nil
}
