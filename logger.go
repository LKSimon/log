package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	DEFAULT_SPLITTYPE           = SplitType_Size
	DEFAULT_FILE_SIZE           = 100
	DEFAULT_FILE_UNIT           = MB
	DEFAULT_CHECK_FILE_INTERNAL = 1 //默认检查文件时间间隔，单位：秒
	DEFAULT_LOG_LEVEL           = TRACE
	DEFAULT_LOG_CHAN_SIZE       = 1000
	DEFAULT_FILE_FORMAT         = TEXT_FORMAT
)

const (
	DATE_FORMAT         = "2006-01-02" //"2006-01-02 15:04:05"
	SUFFIX_FORMAT_DAILY = "2006-01-02" //"2006-01-02 15:04:05"        //按照每日日期切割日志的文件后缀格式
)

//定义单位
type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (10 * iota)
	MB
	GB
	TB
)

//定义文件分离类型
type SplitType byte

const (
	SplitType_Size SplitType = iota
	SplitType_Daily
)

//定义日志层级
type Level byte

//定义json数据类型
type JSONType map[string]interface{}

//定义文件格式
type Format byte

const (
	TEXT_FORMAT Format = iota
	JSON_FORMAT
)

const (
	TRACE Level = iota
	INFO
	WARN
	ERROR
	OFF
)

type FileLogger struct {
	maxFileSize int64
	count       int16     //默认：1,用于按照文件大小切割日志的文件后缀
	flag        int       //默认 log.LstdFlags
	date        time.Time //用于按照日期分割日志
	splitType   SplitType
	level       Level  //默认日志层级：TRACE
	format      Format //定义文件格式
	mutex       *sync.Mutex
	file        *os.File
	lg          *log.Logger
	dir         string //日志存放目录
	name        string //日志文件名
	prefix      string
	logChan     chan string //存放待写入日志信息
}

//默认文件切割类型为:SplitType_Size
func NewDefaultLogger(dir, name string) *FileLogger {
	logger := NewSizeLogger(dir, name, "", DEFAULT_FILE_SIZE, DEFAULT_LOG_CHAN_SIZE, DEFAULT_FILE_UNIT, DEFAULT_FILE_FORMAT)

	return logger
}

func NewSizeLogger(dir, name, prefix string, fileSize, chanSize int64, unit UNIT, format Format) *FileLogger {
	fileName := getFileName(name)
	switch format {
	case TEXT_FORMAT:
		fileName = fileName + ".txt"
	case JSON_FORMAT:
		fileName = fileName + ".json"
	default:
		fileName = fileName + ".txt"
	}

	logger := FileLogger{
		mutex:       new(sync.Mutex),
		dir:         dir,
		name:        fileName,
		level:       DEFAULT_LOG_LEVEL,
		maxFileSize: (fileSize * int64(unit)),
		prefix:      prefix,
		splitType:   SplitType_Size,
		logChan:     make(chan string, chanSize),
		flag:        log.LstdFlags,
		count:       1,
		format:      format,
	}
	logger.initLogger()

	return &logger
}

func NewDailyLogger(dir, name, prefix string, chanSize int64, format Format) *FileLogger {
	fileName := getFileName(name)
	switch format {
	case TEXT_FORMAT:
		fileName = fileName + ".txt"
	case JSON_FORMAT:
		fileName = fileName + ".json"
	default:
		fileName = fileName + ".txt"
	}

	logger := &FileLogger{
		mutex:     new(sync.Mutex),
		dir:       dir,
		name:      name,
		level:     DEFAULT_LOG_LEVEL,
		prefix:    prefix,
		logChan:   make(chan string, chanSize),
		flag:      log.LstdFlags,
		splitType: SplitType_Daily,
	}
	logger.initLogger()

	return logger
}

func (f *FileLogger) initLogger() {
	switch f.splitType {
	case SplitType_Size:
		f.initLoggerBySize()
	case SplitType_Daily:
		f.initLoggerByDaily()
	}
}

func (f *FileLogger) initLoggerBySize() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	logFile := joinFilePath(f.dir, f.name) //生成日志文件绝对路径

	if false == f.isMustSplit() {
		//fmt.Println("不需要切割文件")
		if !isExist(f.dir) { //目录不存在时：
			os.Mkdir(f.dir, 0755)
		}

		f.file, _ = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if f.format != JSON_FORMAT {
			f.lg = log.New(f.file, f.prefix, f.flag)
		}
	} else {
		f.split()
	}

	go f.logWrite()
	go f.fileMonitor()
}

func (f *FileLogger) initLoggerByDaily() {
	f.date, _ = time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))

	f.mutex.Lock()
	defer f.mutex.Unlock()

	logFile := joinFilePath(f.dir, f.name) //生成日志文件绝对路径
	if f.isMustSplit() {                   //f.isMustSplit()已对文件是否存在做出判断，若不存在，则返回false
		f.split()
	} else {
		if !isExist(f.dir) { //文件不存在时：
			os.Mkdir(f.dir, 0755)
		}

		f.file, _ = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		f.lg = log.New(f.file, f.prefix, f.flag)
	}

	go f.logWrite()
	go f.fileMonitor()
}

//判断是否分割文件
func (f *FileLogger) isMustSplit() bool {
	switch f.splitType {
	case SplitType_Size:
		logFile := joinFilePath(f.dir, f.name) //生成日志文件绝对路径
		if !isExist(logFile) {
			return false
		}

		//判断文件大小是否超出初始化时定义的文件大小
		if fileSize(logFile) >= f.maxFileSize {
			return true
		}
	case SplitType_Daily:
		t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
		if t.After(f.date) {
			return true
		}
	}

	return false
}

//定时检查文件大小或日期是否达到切割文件的条件
func (f *FileLogger) fileMonitor() {
	defer func() {
		if r := recover(); nil != r {
			log.Printf("FileLogger'fileMonitor function catch panic: %v\n", r)
		}
	}()

	//定时检查文件是否达到分割文件的条件,暂定一秒检查一次
	ticker := time.NewTicker(time.Duration(DEFAULT_CHECK_FILE_INTERNAL) * time.Second)
	for {
		select {
		case <-ticker.C:
			f.checkFile() //注：f.checkFile()内部已加锁
		}
	}
}

//切割文件
func (f *FileLogger) split() {
	defer func() {
		if r := recover(); nil != r {
			log.Printf("FileLogger'split function catch panic: %v\n", r)
		}
	}()

	logFile := joinFilePath(f.dir, f.name)
	switch f.splitType {
	case SplitType_Size:
		if nil != f.file {
			f.file.Close()
		}

		var bak string
		for i := f.count; ; i++ {
			logFileBak := fmt.Sprint(i, "_", f.name)
			if !isExist(logFileBak) {
				f.count = i
				bak = joinFilePath(f.dir, logFileBak)
				break
			}
		}

		if err := os.Rename(logFile, bak); nil != err {
			//fmt.Println(logFile)
			//fmt.Printf("FileLogger rename error: %v\n", err.Error())
			panic(err.Error())
		} else {
			f.file, _ = os.Create(logFile)
			if f.format != JSON_FORMAT {
				f.lg = log.New(f.file, f.prefix, f.flag)
			}
		}
	case SplitType_Daily:
		logFileBak := fmt.Sprint(time.Now().Format(SUFFIX_FORMAT_DAILY), "_", f.name)
		bak := joinFilePath(f.dir, logFileBak)

		bak = joinFilePath(f.dir, "2019-05-13_"+f.name)
		//fmt.Println("bak: ", bak)
		switch isExist(bak) {
		case false:
			if f.isMustSplit() { //备份文件不存在且需文件分割时：
				if nil != f.file {
					f.file.Close()
				}

				err := os.Rename(logFile, bak)
				if nil != err {
					//fmt.Printf("FileLogger rename error: ", err.Error())
					panic(err.Error())
				} else {
					//fmt.Println("重命名文件成功")
					f.date, _ = time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
					f.file, _ = os.Create(logFile)
					f.lg = log.New(f.file, f.prefix, f.flag)
				}
			}
		case true:
			//fmt.Println("备份文件存在")
		}
	}
}

func (f *FileLogger) checkFile() {
	defer func() {
		if r := recover(); nil != r {
			f.lg.Printf("FileLogger'checkFile function catch panic: %v\n", r)
		}
	}()

	if f.isMustSplit() {
		f.mutex.Lock()
		f.split()
		//fmt.Println("切割完成")
		f.mutex.Unlock()
	}
}

func (f *FileLogger) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	close(f.logChan)
	f.lg = nil

	return f.file.Close()
}
