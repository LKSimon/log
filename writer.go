package log

import (
	"fmt"
	"log"

	//"reflect"
	"runtime"
)

//接收channel 信息并写入文件
func (f *FileLogger) logWrite() {

	defer func() {
		if r := recover(); nil != r {
			log.Printf("FileLogger'logWrite function catch panic: %v\n", r)
		}
	}()

	for {
		select {
		case str := <-f.logChan:
			//fmt.Println("logchan: ", str)
			f.p(str)
		}
	}
}

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

//写入日志
func (f *FileLogger) Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprint(v...)
}

func (f *FileLogger) Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintf(format, v...)
}

func (f *FileLogger) Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	f.logChan <- fmt.Sprintf("%v:%v  ", file, line) + fmt.Sprintln(v...)
}
