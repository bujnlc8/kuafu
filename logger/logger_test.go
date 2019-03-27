package logger_test

import (
	"github.com/linghaihui/kuafu/logger"
	"testing"
)

func TestLoggerFile_Debug(t *testing.T) {
	log := logger.NewLogger("../logs/test.log", 10, 1, logger.LevelWarn)
	log.Debug("testing")
	log.Warn("testing")
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}

	}()
	panic("fatal")
}

func BenchmarkLoggerFile_Debug(b *testing.B) {
	log := logger.NewLogger("../logs/test.log", 0, logger.SplitBySize, logger.LevelDebug)
	for i := 0; i < b.N; i++ {
		if i%4 == 1 {
			log.Error(i)
		} else if i%4 == 2 {
			log.Warn(i)
		} else {
			log.Debug(i)
		}
	}
}
