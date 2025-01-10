package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/sirupsen/logrus"
)

func DingWarn(pattern string, args ...any) error {
	logrus.WithField("args", args).Info(pattern)
	return DingText(alert.SeverityLow, "info", fmt.Sprintf(pattern, args...))
}

func DingWarnf(pattern string, args ...any) error {
	logrus.WithField("args", args).Warn(pattern)
	return DingText(alert.SeverityMedium, "warn", fmt.Sprintf(pattern, args...))
}

func DingError(err error, describe ...string) error {
	if len(describe) == 0 {
		describe = append(describe, "unexpected error")
	}
	logrus.WithError(err).Error(describe)
	return DingText(alert.SeverityHigh, describe[0], fmt.Sprintf("%+v", err))
}

func DingPanicf(err error, description string, args ...any) {
	logrus.WithField("args", args).WithError(err).Error(description)
	err = DingText(alert.SeverityCritical, description, fmt.Sprintf("%+v", err))
	panic(err)
}

func DingText(level alert.Severity, brief, detail string) error {
	ch, ok := alert.DefaultManager().Channel("dingtalk")
	if !ok {
		return errors.New("dingtalk channel not found")
	}

	title := "info"
	switch level {
	case alert.SeverityLow:
		title = "info"
	case alert.SeverityMedium:
		title = "warn"
	case alert.SeverityHigh:
		title = "error"
	case alert.SeverityCritical:
		title = "critical"
	}

	return ch.Send(context.Background(), &alert.Notification{
		Title:    title,
		Content:  fmt.Sprintf("%s\n%s", brief, detail),
		Severity: alert.Severity(level),
	})
}
