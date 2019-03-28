package logger_test

import (
	"github.com/linghaihui/kuafu/logger"
	"os"
	"testing"
)

func TestLoggerSizeFile(t *testing.T) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	log := logger.NewSizeLogger("logs/test-size.log", 0, logger.LevelWarn)
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

func TestLoggerTimeFile(t *testing.T) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	log := logger.NewTimeLogger("logs/test-time.log", logger.EveryDay, logger.LevelWarn)
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

func BenchmarkLoggerSizeFile(b *testing.B) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	log := logger.NewSizeLogger("logs/test-size-bench.log", 0, logger.LevelDebug)
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

func BenchmarkLoggerTimeFile(b *testing.B) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	log := logger.NewTimeLogger("logs/test-time-bench.log", logger.EveryDay, logger.LevelDebug)
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
