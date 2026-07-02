package notify

import (
	"fmt"
	"time"
)

// wechat API 消息结构
type wechatMessage struct {
	MsgType  string          `json:"msgtype"`
	Markdown *wechatMarkdown `json:"markdown,omitempty"`
	Text     *wechatText     `json:"text,omitempty"`
}

type wechatMarkdown struct {
	Content string `json:"content"`
}

type wechatText struct {
	Content string `json:"content"`
}

type wechatResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// defaultBuilder 默认 markdown 告警格式，未注册 bizType 时使用。
var defaultBuilder MessageBuilder = MessageBuilderFunc(buildDefaultMarkdown)

func buildDefaultMarkdown(userId string, bizType, body string) (string, interface{}) {
	now := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(
		"## <font color=\"warning\">服务告警通知</font>\n"+
			">**业务类型：** %s\n"+
			">**用户ID：** %s\n"+
			">**告警内容：** %s\n"+
			">**时间：** %s",
		bizType, userId, body, now,
	)
	return "markdown", &wechatMarkdown{Content: content}
}
