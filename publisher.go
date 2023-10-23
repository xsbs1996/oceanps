package oceanps

import "context"

// Publish 发布事件
func (e *EventTopic) Publish(msg []byte) error {
	err := e.pPManage.PushMsgFn(context.Background(), e.queueName, msg)
	if err != nil {
		return err
	}
	return nil
}
