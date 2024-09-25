package oceanpsfuncs

import (
	"context"
	"errors"
	"sync"
	"time"
)

var PushPullMap = &sync.Map{}

type PushMsgBody struct {
	Msg     []byte        // 消息体
	Exp     time.Duration // 过期时间
	IsAsync bool          // 是否异步
	Key     string
}

type PushPullManage interface {
	CheckClient() error                                                     // 检测链接
	PushMsg(ctx context.Context, topic string, body *PushMsgBody) error     // push消息
	PullMsg(ctx context.Context, topic string, msgChan chan<- []byte) error // pull消息
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
