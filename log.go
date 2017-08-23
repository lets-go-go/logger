package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
)

type LEVEL int
type UNIT int64

type OutputType int
type RollingType int

const (
	ALL LEVEL = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

const (
	Console OutputType = 1 << iota
	File
)

const (
	RollingDaily RollingType = 1 << iota
	RollingSize
)

const (
	_  UNIT = iota
	KB      = 1 << (10 * iota)
	MB
	GB
	TB
)

var (
	// 日志前缀，将写在日期和等级后面，日志内容前面
	prefixes map[LEVEL]string
	// 日志大小单位
	units map[string]UNIT
)

func init() {
	prefixes = map[LEVEL]string{
		TRACE: "[TRACE]",
		DEBUG: "[DEBUG]",
		INFO:  "[INFO ]",
		WARN:  "[WARN ]",
		ERROR: "[ERROR]",
		FATAL: "[FATAL]",
	}

	units = map[string]UNIT{
		"KB": KB,
		"MB": MB,
		"GB": GB,
		"TB": TB,
	}
}

// Config log配置
type Config struct {
	Level                  LEVEL         `json:"Level"`
	OutputType             OutputType    `json:"OutputType"`
	LogFileRollingType     RollingType   `json:"LogFileRollingType"`
	LogFileOutputDir       string        `json:"LogFileOutputDir"`
	LogFileName            string        `json:"LogFileName"`
	LogFileNameDatePattern string        `json:"LogFileNameDatePattern"`
	LogFileNameExt         string        `json:"LogFileNameExt"`
	LogFileMaxCount        int           `json:"LogFileMaxCount"`
	LogFileMaxSize         int64         `json:"LogFileMaxSize"`
	LogFileMaxSizeUnit     string        `json:"LogFileMaxSizeUnit"`
	LogFileScanInterval    time.Duration `json:"LogFileScanInterval"` // 秒
	Sync                   bool          `json:"Sync"`                // 是否同步写
}

// DEFAULT_CONFIG 默认配置
var DEFAULT_CONFIG = &Config{
	Level:                  INFO,
	OutputType:             Console | File,
	LogFileRollingType:     RollingDaily,
	LogFileOutputDir:       ".",
	LogFileName:            "app",
	LogFileNameDatePattern: "20060102",
	LogFileNameExt:         ".log",
	LogFileMaxCount:        5,
	LogFileMaxSize:         5,
	LogFileMaxSizeUnit:     "MB",
	LogFileScanInterval:    1,
	Sync:                   true,
}

type logContent struct {
	level LEVEL
	txt   string
}

// Logger Logger
type Logger struct {
	config *Config
	// 内置logger
	builtInLoggers map[LEVEL]*log.Logger
	// 日志队列
	c chan logContent
	// 当前日志文件
	f *os.File
	// 检查文件monitor是否在运行
	isMonitorRunning bool

	// fileDate 按天rolling的时间
	fileDate string
}

// NewLogger 通过配置项配置
func NewLogger(configStr string) *Logger {
	// 默认配置
	l := &Logger{}
	l.setConfigStr(configStr)
	l.init()
	return l
}

// DefalutConfig 默认配置
func DefalutConfig() *Config {
	return DEFAULT_CONFIG
}

// NewLoggerWithConfig 通过代码配置
func NewLoggerWithConfig(config *Config) *Logger {
	// 默认配置
	l := &Logger{}
	l.setConfig(config)
	l.init()
	return l
}

func (l *Logger) init() {

	flags := log.Ldate | log.Lmicroseconds

	if l.config.Sync {
		flags = flags | log.Lshortfile
	}

	l.builtInLoggers = map[LEVEL]*log.Logger{
		TRACE: log.New(os.Stdout, prefixes[TRACE], flags),
		DEBUG: log.New(os.Stdout, prefixes[DEBUG], flags),
		INFO:  log.New(os.Stdout, prefixes[INFO], flags),
		WARN:  log.New(os.Stdout, prefixes[WARN], flags),
		ERROR: log.New(os.Stdout, prefixes[ERROR], flags),
		FATAL: log.New(os.Stdout, prefixes[FATAL], flags),
	}

	if !l.config.Sync {
		l.c = make(chan logContent, 5000)
		// log write
		go func() {
			for {
				select {
				case content := <-l.c:
					{
						l.builtInLoggers[content.level].Output(3, content.txt)
					}

				}
			}
		}()
	}

}

func (l *Logger) setConfigStr(configStr string) {
	var config Config
	json.Unmarshal([]byte(configStr), &config)

	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		log.Println("=========== parse config failed!!! ==========", err)
		return
	}

	unit := strings.ToUpper(config.LogFileMaxSizeUnit)
	config.LogFileMaxSize = config.LogFileMaxSize * int64(units[unit])

	config.LogFileScanInterval = config.LogFileScanInterval * time.Second
	l.setConfig(&config)
}

func (l *Logger) setConfig(c *Config) {
	l.config = c

	unit := strings.ToUpper(l.config.LogFileMaxSizeUnit)
	l.config.LogFileMaxSize = l.config.LogFileMaxSize * int64(units[unit])
	l.config.LogFileScanInterval = l.config.LogFileScanInterval * time.Second
	l.startFileCheckMonitor()
}

// Output 输出日志
func (l *Logger) Output(level LEVEL, txt string) {

	if fwLogger == nil {
		log.Println("logger not initialed")
		return
	}

	if level >= l.config.Level {
		// l.c <- content
		if l.config.OutputType&File == File {
			if l.f == nil {
				l.makeFile()
			}

			if l.config.Sync {
				l.builtInLoggers[level].Output(3, txt)
			} else {
				l.c <- logContent{level, txt}
			}
		}
	}
}

func (l *Logger) startFileCheckMonitor() {
	if l.isMonitorRunning {
		return
	}
	l.isMonitorRunning = true
	// file check monitor
	go func() {
		monitorTimer := time.NewTicker(l.config.LogFileScanInterval)
		for {
			select {
			case <-monitorTimer.C:
				l.checkFile()
			}
		}
	}()
}

// 初始化日志文件
func (l *Logger) makeFile() {
	if l.config.OutputType == Console {
		return
	}
	if l.f == nil {
		var err error
		var fileName = l.config.LogFileName
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			l.fileDate = time.Now().Format(l.config.LogFileNameDatePattern)
			fileName += "-" + l.fileDate
		}
		if l.config.LogFileRollingType&RollingSize == RollingSize {
			fileName += "-" + l.genFileSeq()
		}

		fullPath := filepath.Join(l.config.LogFileOutputDir, fileName+l.config.LogFileNameExt)

		if l.f, err = os.OpenFile(fullPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
			log.Println("=========== create log file failed!!! ========", err)
		}
	}

	l.updateWriter()

	l.builtInLoggers[INFO].Output(3, "============Welcome to user lets-go-go logger===============")

	l.updateFiles()
}

// updateWriter update writer
func (l *Logger) updateFiles() {

	// remove files
	if l.config.LogFileMaxCount > 0 {
		if list, err := ioutil.ReadDir(l.config.LogFileOutputDir); err == nil {

			var logList []os.FileInfo
			for _, fi := range list {
				if path.Ext(fi.Name()) == l.config.LogFileNameExt {
					logList = append(logList, fi)
				}
			}
			sort.Slice(logList, func(i, j int) bool { return logList[i].ModTime().Before(logList[j].ModTime()) })

			if len(logList) > l.config.LogFileMaxCount {
				for index := 0; index < len(logList)-l.config.LogFileMaxCount; index++ {
					toDeleteFilePath := filepath.Join(l.config.LogFileOutputDir, logList[index].Name())

					os.Remove(toDeleteFilePath)
					txt := fmt.Sprintf("remove older log files:%v, modtime=%v", toDeleteFilePath, logList[index].ModTime())
					l.builtInLoggers[INFO].Output(3, txt)
				}
			}
		}
	}
}

// updateWriter update writer
func (l *Logger) updateWriter() {
	logWriters := []io.Writer{l.f}
	if l.config.OutputType&Console == Console {

		logWriters = append(logWriters, os.Stdout)
	}

	for _, builtInLogger := range l.builtInLoggers {
		builtInLogger.SetOutput(io.MultiWriter(logWriters...))
	}
}

// 检查文件是否需要重新创建
func (l *Logger) checkFile() {
	if l.config.OutputType == Console || l.f == nil {
		return
	}
	needRecreate := false
	if l.config.LogFileRollingType&RollingDaily == RollingDaily {
		currentDate := time.Now().Format(l.config.LogFileNameDatePattern)
		if currentDate != l.fileDate {
			needRecreate = true
		}
	} else if l.config.LogFileRollingType&RollingSize == RollingSize {
		info, err := os.Stat(filepath.Join(l.config.LogFileOutputDir, l.f.Name()))
		if err != nil {
			log.Println("============= check file size failed!!! ==========", err)
			return
		}
		if info.Size() >= l.config.LogFileMaxSize {
			log.Println("============= need recreate ==========", info.Size())
			needRecreate = true
		}
	}

	if needRecreate {
		l.f.Close()
		l.f = nil

		l.makeFile()
	}
}

// 生成日志文件序列号，并保存到.seq
func (l *Logger) genFileSeq() string {
	seqFile := filepath.Join(l.config.LogFileOutputDir, ".seq")
	if IsFileExists(seqFile) {
		if bytes, err := ioutil.ReadFile(seqFile); err == nil {
			if seq, err := strconv.Atoi(string(bytes)); err == nil {
				ioutil.WriteFile(seqFile, []byte(strconv.Itoa(seq+1)), 0666)
				return strconv.Itoa(seq + 1)
			}
		}
	}
	ioutil.WriteFile(seqFile, []byte("1"), 0666)
	return "1"
}

// 重置日志文件序列号
func (l *Logger) resetFileSeq() {
	seqFile := filepath.Join(l.config.LogFileOutputDir, ".seq")
	if IsFileExists(seqFile) {
		ioutil.WriteFile(seqFile, []byte("1"), 0666)
	}
}

// IsFileExists 判断文件是否存在
func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
