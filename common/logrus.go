package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogrus() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true) // 打印调用者信息
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 处理函数名（保持原样）
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			// 处理文件名：计算相对路径
			// 获取当前工作目录 (cwd)
			cwd, _ := os.Getwd()
			// 计算 f.File 相对于 cwd 的路径
			relPath, err := filepath.Rel(cwd, f.File)
			if err != nil {
				// 如果计算失败（比如文件在不同磁盘），回退显示完整路径或文件名
				relPath = f.File
			}
			return funcname, fmt.Sprintf("%s:%d", relPath, f.Line)
		},
	})
	logger := &lumberjack.Logger{
		Filename:   "mq.log",
		MaxSize:    1,    // 日志文件最大大小（MB）
		MaxBackups: 10,   // 保留的旧日志文件数量
		MaxAge:     30,   // 保留的旧日志文件天数
		Compress:   true, // 是否压缩旧日志文件
	}

	// 设置日志输出到 lumberjack
	logrus.SetOutput(logger)
}
