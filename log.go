/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
 Package: log
*/
package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const DateFormat = "2006-01-02"
const TimeFormat = "2006-01-02 15:04:05"

type LEVEL byte

const (
	TRACE LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	OFF
)

type LoggerConf struct {
	FileDir  string
	FileName string
	Prefix   string
	Level    string
}

var (
	fileDir  string
	fileName string
	prefix   string
	date     *time.Time
	logFile  *os.File
	logger   *log.Logger
	logLevel LEVEL
	mutex    *sync.RWMutex
	logChan  chan string
)

// 初始化日志配置
func BootLogger() (err error) {
	conf := &LoggerConf{
		FileDir:  GetLogsDir(),
		FileName: GetLogsFilename(),
		Prefix:   GetLogsPrefix(),
		Level:    GetLogsLevel(),
	}

	fileDir = conf.FileDir
	fileName = conf.FileName
	prefix = conf.Prefix
	mutex = new(sync.RWMutex)
	logChan = make(chan string, 8000)

	if strings.EqualFold(conf.Level, "OFF") {
		logLevel = OFF
	} else if strings.EqualFold(conf.Level, "TRACE") {
		logLevel = TRACE
	} else if strings.EqualFold(conf.Level, "INFO") {
		logLevel = INFO
	} else if strings.EqualFold(conf.Level, "WARN") {
		logLevel = WARN
	} else if strings.EqualFold(conf.Level, "ERROR") {
		logLevel = ERROR
	} else {
		logLevel = DEBUG
	}

	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	date = &t

	if isMustSplit() {
		if err = split(); err != nil {
			return
		}

	} else {
		isExistOrCreate()

		logFilepath := filepath.Join(fileDir, fileName)
		logFile, err = os.OpenFile(logFilepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}

		logger = log.New(logFile, prefix, log.LstdFlags|log.Lmicroseconds)
	}

	go logWriter()
	go fileMonitor()

	return
}

// 日志文件是否分割
func isMustSplit() bool {
	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	return t.After(*date)
}

// 检查日志文件目录是否存在，不存在则创建
func isExistOrCreate() {
	_, err := os.Stat(fileDir)
	if err != nil && !os.IsExist(err) {
		mkdirErr := os.Mkdir(fileDir, 0755)
		if mkdirErr != nil {
			log.Println("Create dir failed, error: ", mkdirErr)
		}
	}
}

// 分割日志
func split() (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	sourceLog := filepath.Join(fileDir, fileName)
	targetLog := sourceLog + "." + date.Format(DateFormat)

	if logFile != nil {
		_ = logFile.Close()
	}

	err = os.Rename(sourceLog, targetLog)
	if err != nil {
		return
	}

	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	date = &t

	logFile, err = os.OpenFile(sourceLog, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	logger = log.New(logFile, prefix, log.LstdFlags|log.Lmicroseconds)
	return
}

// 日志写入
func logWriter() {
	defer func() { recover() }()

	for {
		str := <-logChan
		mutex.RLock()
		_ = logger.Output(2, str)
		mutex.RUnlock()
	}
}

// 日志分割监控
func fileMonitor() {
	defer func() { recover() }()

	timer := time.NewTicker(30 * time.Second)
	for {
		<-timer.C

		if isMustSplit() {
			if err := split(); err != nil {
				Error("Log split error: %v\n", err)
			}
		}
	}
}

// 关闭日志
func CloseLogger() {
	if logChan != nil {
		close(logChan)
		logger = nil
		_ = logFile.Close()
	}
}

// 输出格式化日志
func Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logChan <- fmt.Sprintf("[%v:%v]", fmt.Sprintf(format, v...)+filepath.Base(file), line)
}

// 输出格式化日志
func Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logChan <- fmt.Sprintf("[%v:%v]", fmt.Sprint(v...)+filepath.Base(file), line)
}

// 输出格式化日志
func Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintln(v...)
}

// 输出致命错误日志, 并退出系统
func Fatal(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line) + fmt.Sprintln(v...)
	_ = log.Output(2, fmt.Sprintln(v))
	os.Exit(1)
}

// 输出致命错误日志, 并退出系统
func Fatally(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line) + fmt.Sprintln(v...)
	_ = log.Output(2, fmt.Sprintln(v))
	os.Exit(1)
}

// 输出跟踪日志
func Trace(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	if logLevel <= TRACE {
		logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[TRACE] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出调试日志
func Debug(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%v:%v:%v%v]", fmt.Sprintf("[DEBUG] [")+filepath.Base(file), line, format, v)
	fmt.Printf("%s\033[0;40;34m%s\033[0m\n", setNowTime(), s)
	if logLevel <= DEBUG {
		logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[DEBUG] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出信息日志
func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%v:%v:%v%v]", fmt.Sprintf("[INFO] [")+filepath.Base(file), line, format, v)
	fmt.Printf("%s\033[0;40;32m%s\033[0m\n", setNowTime(), s)
	if logLevel <= INFO {
		logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[INFO] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出警告日志
func Warning(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%v:%v:%v%v]", fmt.Sprintf("[WARN] [")+filepath.Base(file), line, format, v)
	fmt.Printf("%s\033[0;40;33m%s\033[0m\n", setNowTime(), s)
	if logLevel <= WARN {
		logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[WARN] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出错误日志
func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%v:%v:%v%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line, format, v)
	fmt.Printf("%s\033[0;40;31m%s\033[0m\n", setNowTime(), s)
	if logLevel <= ERROR {
		logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出格式化后的当前时间字符串
func setNowTime() string {
	return time.Now().Format(TimeFormat)
}
