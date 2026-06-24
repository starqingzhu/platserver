package notify

import (
	"fmt"
	"golib/zaplog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.yaml.in/yaml/v3"
)

const configFileName = "RobotNotice.yaml"

type robotNoticeConfig struct {
	WebhookURL string `yaml:"webhookURL"`
}

func LoadConfig(configDir string) error {
	filePath := filepath.Join(configDir, configFileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	var cfg robotNoticeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse %s: %w", configFileName, err)
	}

	if cfg.WebhookURL == "" {
		return fmt.Errorf("%s: webhookURL is empty", configFileName)
	}

	n := &RobotNotifier{
		webhookURL: cfg.WebhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
		queue:      make(chan sendTask, queueSize),
	}
	n.startWorker()
	defaultNotifier = n
	zaplog.LoggerSugar.Infof("[notify] robot notifier loaded from %s", configFileName)
	return nil
}
