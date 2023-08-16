package utils

import (
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/sirupsen/logrus"
)

func DingWarn(pattern string, args ...any) {
	logrus.WithField("args", args).Info(pattern)
	DingText("ℹ️ info", "info", fmt.Sprintf(pattern, args...))
}

func DingWarnf(pattern string, args ...any) {
	logrus.WithField("args", args).Warn(pattern)
	DingText("⚠️ warn", "warn", fmt.Sprintf(pattern, args...))
}

func DingError(err error, describe ...string) {
	if len(describe) == 0 {
		describe = append(describe, "unexpected error")
	}
	logrus.WithError(err).Error(describe)
	DingText("💔 error", describe[0], fmt.Sprintf("%+v", err))
}

func DingPanicf(err error, description string, args ...any) {
	logrus.WithField("args", args).WithError(err).Error(description)
	DingText("😱 panic", description, fmt.Sprintf("%+v", err))
	panic(err)
}

func DingText(level, brief, detail string) {
	if err := alert.SendDingTalkTextMessage(level, brief, detail); err != nil {
		logrus.WithError(err).WithField("detail", detail).Error("failed to send dingding")
	}
}
