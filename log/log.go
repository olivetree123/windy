package log

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

func init() {
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// Info 提示信息
//func Info(args ...interface{}) {
//	logger.Info(args)
//}

// Warn 警告信息
//func Warn(args ...interface{}) {
//	logger.Warn(args)
//}

// Error 错误信息
//func Error(args ...interface{}) {
//	logger.Error(args)
//}
