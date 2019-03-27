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
	LevelDebug = int8(1)
	LevelWarn  = int8(2)
	LevelErr   = int8(3)
	LevelFatal = int8(4)
)

type LoggerFile struct {
	LoggerAbstract
	lock        *sync.RWMutex
	path        string
	maxSize     int64
	file        *os.File
	splitByDate bool
	splitBySize bool
	minLevel    int8
}

func NewLogger(path string, maxSize int64, splitByWhat int, minLevel int8) *LoggerFile {
	if fd, err := os.Create(path); err != nil {
		panic(util.FormatString("open path %s error happens", path))
	} else {
		var splitByDate, splitBySize bool
		if splitByWhat == SplitBySize {
			splitBySize = true
		} else if splitByWhat == SplitByDate {
			splitByDate = true
		}
		if maxSize == 0 {
			maxSize = 1024 * 1024 * 2 //default is 2MB
		}
		if minLevel == 0 {
			minLevel = LevelWarn
		}
		return &LoggerFile{
			lock:        new(sync.RWMutex),
			path:        path,
			maxSize:     maxSize,
			file:        fd,
			splitByDate: splitByDate,
			splitBySize: splitBySize,
			minLevel:    minLevel,
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
	if logger.splitBySize {
		if fileInfo, err := os.Stat(logger.path); err != nil {
			panic(err)
		} else {
			if fileInfo.Size() > logger.maxSize {
				if err := os.Rename(logger.path, logger.path+"."+strconv.Itoa(time.Now().Nanosecond())); err != nil {
					panic(err)
				}
				if newFile, err := os.Create(logger.path); err != nil {
					panic(err)
				} else {
					logger.file = newFile
				}
			}
		}
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
