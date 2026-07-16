package notify

import (
	"fmt"
	"strings"
	"time"
)

func init() {
	RegisterBuilder("activity", MessageBuilderFunc(buildActivityMarkdown))
}

func buildActivityMarkdown(userId string, bizType, body string) (string, interface{}) {
	now := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")

	// WeChat markdown supports: info(gray) comment(gray) warning(yellow) red(red)
	color := "warning" // default yellow for modify and other operations
	firstLine := body
	if idx := strings.Index(body, "\n"); idx >= 0 {
		firstLine = body[:idx]
	}
	if strings.Contains(strings.ToLower(firstLine), "[delete]") {
		color = "red"
	}

	// Format each body line with double blockquote for indented appearance
	lines := strings.Split(body, "\n")
	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(">>")
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	bodyBlock := strings.TrimRight(sb.String(), "\n")

	content := fmt.Sprintf(
		"## <font color=\"%s\">服务告警通知</font>\n"+
			">**业务类型：** %s\n"+
			">**用户ID：** %s\n"+
			">**告警内容：**\n"+
			"%s\n"+
			">**时间：** %s",
		color, bizType, userId, bodyBlock, now,
	)
	return "markdown", &wechatMarkdown{Content: content}
}
