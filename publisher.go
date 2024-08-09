package oceanps

import (
	"context"
	"time"
)

// Publish 发布事件
func (e *EventTopic) Publish(ctx context.Context, msg []byte, exp time.Duration) error {
	err := e.pPManage.PushMsgFn(ctx, e.queueName, msg, exp)
	if err != nil {
		return err
	}
	return nil
}
