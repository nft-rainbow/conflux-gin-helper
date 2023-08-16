package logger

import (
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	pprofLogOutput io.Writer
	logConfig      LogConfig
)

type LogConfig struct {
	Level  string `yaml:"level"`
	Folder string `yaml:"folder"`
}

func Init(config LogConfig, logo string) {
	logConfig = config

	log.SetOutput(output())
	log.Print(logo)

	GetPprofLogOutput()

	logrus.SetOutput(output())
	logrus.SetLevel(level())
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "ztimestamp",
			logrus.FieldKeyLevel: "zlevel",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyFunc:  "zcaller",
			logrus.FieldKeyFile:  "zfile",
		},
	})

	// logrus.AddHook(&LogHook{})
	logrus.Info("init logrus done")
}

func RefreshOutput() {
	logrus.SetOutput(output())
}

func GetPprofLogOutput() io.Writer {
	if pprofLogOutput != nil {
		return pprofLogOutput
	}
	slowLogFile := path.Join(logConfig.Folder, "slow.log")
	f, err := os.OpenFile(slowLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	pprofLogOutput = f
	return pprofLogOutput
}

func output() io.Writer {
	dirPath := logConfig.Folder
	if dirPath == "" {
		logrus.Panic("not set log folder path")
	}
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			logrus.Panic(err)
		}
	}

	filePath := path.Join(dirPath, time.Now().Format("2006_01_02.log"))
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Panic(err)
	}
	return io.MultiWriter(os.Stdout, writer)
}

func level() logrus.Level {
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	lvl := logConfig.Level
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.Panic(err)
	}
	return level
}
