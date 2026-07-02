package notify

import "fmt"

// MessageBuilder 根据业务参数构建发送给企业微信的消息体。
type MessageBuilder interface {
	Build(userId string, bizType, body string) (msgType string, payload interface{})
}

// MessageBuilderFunc 函数类型实现 MessageBuilder，方便轻量注册。
type MessageBuilderFunc func(userId string, bizType, body string) (msgType string, payload interface{})

func (f MessageBuilderFunc) Build(userId string, bizType, body string) (string, interface{}) {
	return f(userId, bizType, body)
}

var builderRegistry = map[string]MessageBuilder{}

// RegisterBuilder 注册指定 bizType 的消息构建器，在 init() 中调用。
func RegisterBuilder(bizType string, b MessageBuilder) {
	if _, exists := builderRegistry[bizType]; exists {
		panic(fmt.Sprintf("[notify] builder already registered for bizType: %s", bizType))
	}
	builderRegistry[bizType] = b
}

func getBuilder(bizType string) MessageBuilder {
	if b, ok := builderRegistry[bizType]; ok {
		return b
	}
	return defaultBuilder
}
