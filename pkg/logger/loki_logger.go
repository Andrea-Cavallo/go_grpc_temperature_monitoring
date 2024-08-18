package logger

import (
	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/loki-client-go/pkg/urlutil"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	BatchWait = 1 * time.Second
	BatchSize = 1024 * 10
)

// LokiHook is a custom Logrus hook that sends logs to Loki
type LokiHook struct {
	client *loki.Client
	levels []logrus.Level
}

// NewLokiHook creates a new LokiHook
func NewLokiHook(lokiURL string, levels []logrus.Level) (*LokiHook, error) {
	var parsedURL urlutil.URLValue
	if err := parsedURL.Set(lokiURL); err != nil {
		return nil, errors.Wrap(err, "failed to parse Loki URL")
	}
	cfg := loki.Config{
		URL:       parsedURL,
		BatchWait: BatchWait,
		BatchSize: BatchSize,
		Timeout:   30 * time.Second,
	}

	client, err := loki.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Loki client")
	}

	return &LokiHook{
		client: client,
		levels: levels,
	}, nil
}

// Levels returns the log levels supported by this hook
func (hook *LokiHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire is called by Logrus when a log entry needs to be sent to Loki
func (hook *LokiHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	labels := createLabels(entry.Level.String(), entry.Data)

	err = hook.client.Handle(
		labels,
		entry.Time,
		line,
	)

	if err != nil {
		// Log the error locally
		logrus.Errorf("Failed to send log to Loki: %v", err)
		return err
	}

	return nil
}

// createLabels creates labels for Loki log messages
func createLabels(level string, data logrus.Fields) model.LabelSet {
	labels := model.LabelSet{
		"level": model.LabelValue(level),
	}

	// Add additional fields as labels if needed
	for key, value := range data {
		labels[model.LabelName(key)] = model.LabelValue(value.(string))
	}

	return labels
}

// Close closes the Loki client and ensures all logs are sent
func (hook *LokiHook) Close() {
	hook.client.Stop()
}
