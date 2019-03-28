package logger

import (
	"fmt"
	"github.com/linghaihui/kuafu/util"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SplitByDate = 1
	SplitBySize = 2
)

const (
	EveryHour  = 1
	HalfDay    = 2
	EveryDay   = 3
	EveryWeek  = 4
	EveryMonth = 5
)

const (
	LevelDebug = 1
	LevelWarn  = 2
	LevelErr   = 3
	LevelFatal = 4
)

type LoggerFile struct {
	LoggerAbstract
	lock          *sync.RWMutex
	path          string
	maxSize       int64
	file          *os.File
	splitBy       int
	minLevel      int8
	splitTimeMode int
	nextTime      time.Time
}

// new split by size logger
func NewSizeLogger(path string, maxSize int64, minLevel int8) *LoggerFile {
	if fd, err := os.Create(path); err != nil {
		panic(util.FormatString("open path %s error happens", path))
	} else {
		if maxSize == 0 {
			maxSize = 1024 * 1024 * 2 //default is 2MB
		}
		if minLevel == 0 {
			minLevel = LevelWarn // default log level is warn
		}
		return &LoggerFile{
			lock:     new(sync.RWMutex),
			path:     path,
			maxSize:  maxSize,
			file:     fd,
			splitBy:  SplitBySize,
			minLevel: minLevel,
		}
	}
	return nil
}

// gen the next time of specific splitTime mode
func genNextTime(splitTimeMode int) time.Time {
	switch splitTimeMode {
	case EveryHour:
		nowStr := time.Now().Format("2006-01-02 15:00:00")
		now, _ := time.Parse("2006-01-02 15:00:00", nowStr)
		now = now.Add(time.Duration(time.Hour))
		return now
	case HalfDay:
		now := time.Now()
		if now.Hour() >= 12 {
			nowStr := time.Now().Format("2006-01-02 00:00:00")
			now, _ = time.Parse("2006-01-02 00:00:00", nowStr)
			now = now.Add(time.Duration(24 * time.Hour))
		} else {
			nowStr := time.Now().Format("2006-01-02 00:00:00")
			now, _ = time.Parse("2006-01-02 00:00:00", nowStr)
			now = now.Add(time.Duration(12 * time.Hour))
		}
		return now
	case EveryDay:
		nowStr := time.Now().Format("2006-01-02 00:00:00")
		now, _ := time.Parse("2006-01-02 00:00:00", nowStr)
		now = now.Add(time.Duration(24 * time.Hour))
		return now
	case EveryWeek:
		nowStr := time.Now().Format("2006-01-02 00:00:00")
		now, _ := time.Parse("2006-01-02 00:00:00", nowStr)
		noWeekDay := now.Weekday()
		if noWeekDay == time.Sunday {
			now = now.Add(time.Duration(24 * time.Hour))
		} else {
			now = now.Add(time.Duration(24*(8-int(noWeekDay))) * time.Hour)
		}
		return now
	case EveryMonth:
		nowStr := time.Now().Format("2006-01-02 00:00:00")
		now, _ := time.Parse("2006-01-02 00:00:00", nowStr)
		now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0)
		return now
	default:
		panic("do not support split time mode")
	}
	return time.Time{}
}

// new split by time logger
func NewTimeLogger(path string, splitTimeMode int, minLevel int8) *LoggerFile {
	if fd, err := os.Create(path); err != nil {
		panic(util.FormatString("open path %s error happens", path))
	} else {
		if splitTimeMode == 0 {
			splitTimeMode = EveryDay //default is every day
		}
		if minLevel == 0 {
			minLevel = LevelWarn // default log level is warn
		}
		return &LoggerFile{
			lock:          new(sync.RWMutex),
			path:          path,
			file:          fd,
			splitBy:       SplitByDate,
			minLevel:      minLevel,
			splitTimeMode: splitTimeMode,
			nextTime:      genNextTime(splitTimeMode),
		}
	}
	return nil
}

// debug will just log the info
func (logger *LoggerFile) Debug(args ...interface{}) {
	newArgs := []interface{}{"[DEBUG]", time.Now().Format("2006-01-02 15:04:05")}
	newArgs = append(newArgs, args...)
	if logger.minLevel <= LevelDebug {
		logger.write(newArgs...)
	} else {
		fmt.Println(newArgs...)
	}
}

// fatal will just log the info
func (logger *LoggerFile) Warn(args ...interface{}) {
	newArgs := []interface{}{"[WARN]", time.Now().Format("2006-01-02 15:04:05")}
	newArgs = append(newArgs, args...)
	if logger.minLevel <= LevelWarn {
		logger.write(newArgs...)
	} else {
		fmt.Println(newArgs...)
	}
}

// error will log the info and print stacktrace
func (logger *LoggerFile) Error(args ...interface{}) {
	newArgs := []interface{}{"[ERR]", time.Now().Format("2006-01-02 15:04:05")}
	newArgs = append(newArgs, args...)
	buff := stackTrace(true)
	if logger.minLevel <= LevelErr {
		newArgs = append(newArgs, "\nstackTrace:\n"+string(buff))
		logger.write(newArgs...)
	} else {
		fmt.Println(newArgs...)
	}
}

// fatal will log the info and print stacktrace , most importantly , it will exit
func (logger *LoggerFile) Fatal(args ...interface{}) {
	newArgs := []interface{}{"[Fatal]", time.Now().Format("2006-01-02 15:04:05")}
	newArgs = append(newArgs, args...)
	buff := stackTrace(true)
	if logger.minLevel <= LevelFatal {
		newArgs = append(newArgs, "\nstackTrace:\n"+string(buff))
		logger.write(newArgs...)
	} else {
		fmt.Println(newArgs...)
	}
	os.Exit(1)
}
func (logger *LoggerFile) write(args ...interface{}) {
	str := util.FormatString(strings.Repeat("%v  ", len(args)), args...)
	logger.lock.Lock()
	defer logger.lock.Unlock()
	str += "\n"
	switch logger.splitBy {
	case SplitBySize:
		if fileInfo, err := os.Stat(logger.path); err != nil {
			panic(err)
		} else {
			if fileInfo.Size() > logger.maxSize {
				if err := os.Rename(logger.path,
					logger.path+"."+strconv.Itoa(time.Now().Nanosecond())); err != nil {
					panic(err)
				}
				if newFile, err := os.Create(logger.path); err != nil {
					panic(err)
				} else {
					logger.file = newFile
				}
			}
		}
	case SplitByDate:
		if time.Now().After(logger.nextTime) {
			if err := os.Rename(logger.path,
				logger.path+"."+strconv.Itoa(time.Now().Nanosecond())); err != nil {
				panic(err)
			}
			if newFile, err := os.Create(logger.path); err != nil {
				panic(err)
			} else {
				logger.file = newFile
				logger.nextTime = genNextTime(logger.splitTimeMode)
			}
		}
	default:
		panic("do not support split type")
	}
	if _, err := logger.file.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func stackTrace(all bool) []byte {
	buff := make([]byte, 10240)
	for {
		size := runtime.Stack(buff, all)
		if size == len(buff) {
			buff = make([]byte, len(buff)<<1)
			continue
		} else {
			buff = buff[:size]
		}
		break
	}
	return buff
}

type LoggerAbstract interface {
	Debug(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}
