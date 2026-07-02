package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golib/zaplog"
	"io"
	"net/http"
)

const queueSize = 256

type sendTask struct {
	userId  string
	bizType string
	body    string
}

type RobotNotifier struct {
	webhookURL string
	client     *http.Client
	queue      chan sendTask
	done       chan struct{}
}

var defaultNotifier *RobotNotifier

// Send enqueues a message to be delivered asynchronously.
// Returns an error only if the notifier is uninitialised or the queue is full.
func Send(userId string, bizType, body string) error {
	if defaultNotifier == nil {
		return fmt.Errorf("notifier not initialised, call LoadConfig first")
	}
	select {
	case defaultNotifier.queue <- sendTask{userId, bizType, body}:
		return nil
	default:
		return fmt.Errorf("notify queue full, message dropped")
	}
}

// Shutdown drains and closes the queue, waiting for the worker to finish all
// pending messages before returning.
func Shutdown() {
	if defaultNotifier != nil {
		defaultNotifier.shutdown()
	}
}

func (n *RobotNotifier) startWorker() {
	n.done = make(chan struct{})
	go func() {
		defer close(n.done)
		for t := range n.queue {
			if err := n.send(t.userId, t.bizType, t.body); err != nil {
				zaplog.LoggerSugar.Warnf("[notify] send failed: %v", err)
			}
		}
	}()
}

func (n *RobotNotifier) shutdown() {
	close(n.queue)
	<-n.done
}

func (n *RobotNotifier) send(userId string, bizType, body string) error {
	msgType, payload := getBuilder(bizType).Build(userId, bizType, body)

	msg := map[string]interface{}{
		"msgtype": msgType,
		msgType:   payload,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(msg); err != nil {
		return fmt.Errorf("marshal wechat message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.webhookURL, &buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var wr wechatResponse
	if err := json.Unmarshal(data, &wr); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	if wr.ErrCode != 0 {
		return fmt.Errorf("wechat api error %d: %s", wr.ErrCode, wr.ErrMsg)
	}
	return nil
}
