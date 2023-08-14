package localutil

import "github.com/sirupsen/logrus"

var (
	UserRegisterLog *logrus.Logger
)

func LoggerInit() {
	// 创建一个新的 Logrus 日志记录器
	UserRegisterLog = logrus.New()
	// 设置日志级别为 Info，这将只显示 Info 及以上级别的日志
	UserRegisterLog.SetLevel(logrus.InfoLevel)
}
