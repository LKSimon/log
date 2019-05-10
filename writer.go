package log

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"
)

const DEFAULT_TIME_FROMAT = "2006-01-02 15:04:05"

type Fields map[string]interface{}
type JsonMap_t struct { //对外暴漏JsonMap_t结构体，用于解析日志
	Time     string `json:"time"` //时间格式DEFAULT_TIME_FROMAT: 2006-01-02 15:04:05
	FileLine string `json:"fileLine"`
	Level    string `json:"level,omitempty"`
	Data     Fields `json:"data"`
}

//接收channel 信息并写入文件
func (f *FileLogger) logWrite() {
	defer func() {
		if r := recover(); nil != r {
			log.Printf("FileLogger'logWrite function catch panic: %v\n", r)
		}
	}()

	//写入txt日志

	switch f.format {
	case JSON_FORMAT:
		fmt.Println("写入json数据")
		for {
			select {
			case str := <-f.logChan:
				//fmt.Println("logchan: ", str)
				f.json(str)
			}
		}
	case TEXT_FORMAT:
		for {
			select {
			case str := <-f.logChan:
				//fmt.Println("logchan: ", str)
				f.p(str)
			}
		}
	default:
		for {
			select {
			case str := <-f.logChan:
				//fmt.Println("logchan: ", str)
				f.p(str)
			}
		}
	}

	/*
		for {
			select {
			case str := <-f.logChan:
				//fmt.Println("logchan: ", str)
				f.p(str)
			}
		}
	*/
}

//写入json格式数据
func (f *FileLogger) json(str string) { //str为json序列号为map后的字符串
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if nil == f.file {
		fmt.Println("f.file is nil")
	} else {
		f.file.WriteString(str)
		f.file.Sync()
	}
}

//写入txt格式数据
func (f *FileLogger) p(str string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if nil == f.lg {
		fmt.Println("f.lg is nil")
	} else {
		//fmt.Println("f.lg is not nil")
		f.lg.Output(2, str)
	}
}

//写入txt格式日志
func (f *FileLogger) Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprint(v...)
}

//写入txt格式日志
func (f *FileLogger) Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf(format, v...)
}

func (f *FileLogger) Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintln(v...)
}

//日志信息全部写入同一文件时，可使用Tracef、Infof、Warnf、Errorf方法
//trace日志
func (f *FileLogger) Tracef(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if f.level <= TRACE {
		f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf("\033[32m[TRACE]"+format+"\033[0m ", v...)
	}
}

//info日志
func (f *FileLogger) Infof(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if f.level <= INFO {
		f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf("\033[1;35m[INFO]"+format+"\033[0m ", v...)
	}
}

//warn日志
func (f *FileLogger) Warnf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if f.level <= WARN {
		f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf("\033[1;33m[WARN]"+format+"\033[0m ", v...)
	}
}

//error 日志
func (f *FileLogger) Errorf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if f.level <= ERROR {
		f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf("\033[1;31m[ERROR] "+format+" \033[0m ", v...)
	}
}

//写入json格式日志
func (f *FileLogger) PrintJson(fields Fields) {
	defer func() {
		if r := recover(); nil != r {
			log.Println(r)
		}
	}()

	_, file, line, _ := runtime.Caller(1)
	fileLine := fmt.Sprint(file, ": ", line)
	t := time.Now().Format(DEFAULT_TIME_FROMAT)

	ds := JsonMap_t{ //data struct
		FileLine: fileLine,
		Time:     t,
		Data:     fields,
	}

	b, err := json.Marshal(ds)
	if nil != err {
		panic(err)
	}
	b = append(b, '\n')

	f.logChan <- string(b)
}

//日志信息全部写入同一文件时，可使用Tracef、Infof、Warnf、Errorf方法
//trace日志
func (f *FileLogger) TraceJson(fields Fields) {
	if f.level <= TRACE {
		defer func() {
			if r := recover(); nil != r {
				log.Println(r)
			}
		}()

		_, file, line, _ := runtime.Caller(1)
		fileLine := fmt.Sprint(file, ": ", line)
		t := time.Now().Format(DEFAULT_TIME_FROMAT)

		ds := JsonMap_t{ //data struct
			FileLine: fileLine,
			Time:     t,
			Level:    "TRACE",
			Data:     fields,
		}

		b, err := json.Marshal(ds)
		if nil != err {
			panic(err)
		}
		b = append(b, '\n')

		f.logChan <- string(b)
	}
}

//info日志
func (f *FileLogger) InfoJson(fields Fields) {
	if f.level <= INFO {
		defer func() {
			if r := recover(); nil != r {
				log.Println(r)
			}
		}()

		_, file, line, _ := runtime.Caller(1)
		fileLine := fmt.Sprint(file, ": ", line)
		t := time.Now().Format(DEFAULT_TIME_FROMAT)

		ds := JsonMap_t{ //data struct
			FileLine: fileLine,
			Time:     t,
			Level:    "INFO",
			Data:     fields,
		}

		b, err := json.Marshal(ds)
		if nil != err {
			panic(err)
		}
		b = append(b, '\n')

		f.logChan <- string(b)
	}
}

//warn日志
func (f *FileLogger) WarnJson(fields Fields) {
	if f.level <= WARN {
		defer func() {
			if r := recover(); nil != r {
				log.Println(r)
			}
		}()

		_, file, line, _ := runtime.Caller(1)
		fileLine := fmt.Sprint(file, ": ", line)
		t := time.Now().Format(DEFAULT_TIME_FROMAT)

		ds := JsonMap_t{ //data struct
			FileLine: fileLine,
			Time:     t,
			Level:    "WARN",
			Data:     fields,
		}

		b, err := json.Marshal(ds)
		if nil != err {
			panic(err)
		}
		b = append(b, '\n')

		f.logChan <- string(b)
	}
}

//error 日志
func (f *FileLogger) ErrorJson(fields Fields) {
	if f.level <= ERROR {
		defer func() {
			if r := recover(); nil != r {
				log.Println(r)
			}
		}()

		_, file, line, _ := runtime.Caller(1)
		fileLine := fmt.Sprint(file, ": ", line)
		t := time.Now().Format(DEFAULT_TIME_FROMAT)

		ds := JsonMap_t{ //data struct
			FileLine: fileLine,
			Time:     t,
			Level:    "ERROR",
			Data:     fields,
		}

		b, err := json.Marshal(ds)
		if nil != err {
			panic(err)
		}
		b = append(b, '\n')

		f.logChan <- string(b)
	}
}
