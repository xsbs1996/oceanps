package oceanpsfuncs

import (
	"context"
	"errors"
	"sync"
)

var PushPullMap = &sync.Map{}

type PushPullManage interface {
	CheckClient() error                                                       // 检测链接
	PushMsgFn(ctx context.Context, topic string, msg []byte) error            // push消息
	PullMsgFn(ctx context.Context, topic string, msgChan chan<- []byte) error // pull消息
}

// RegisterPushPull 注册
func RegisterPushPull(method string, pushPull PushPullManage) {
	PushPullMap.Store(method, pushPull)
}

// GetPushPull 获取
func GetPushPull(method string) (PushPullManage, error) {
	p, ok := PushPullMap.Load(method)
	if !ok {
		return nil, errors.New("method don't exist")
	}
	return p.(PushPullManage), nil
}
