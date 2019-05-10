package log

//设置日志flag，即标准库日志flag
/*
const (
    Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
    Ltime                         // the time in the local time zone: 01:23:23
    Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
    Llongfile                     // full file name and line number: /a/b/c/d.go:23
    Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
    LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
    LstdFlags     = Ldate | Ltime // initial values for the standard logger
)
*/

func (f *FileLogger) SetFlag(flag int) {
	f.mutex.Lock()
	f.flag = flag
	f.mutex.Unlock()
}

//设置日志层级
func (f *FileLogger) SetLevel(level Level) {
	f.mutex.Lock()
	f.level = level
	f.mutex.Unlock()
}

//设置日志prefix
func (f *FileLogger) Setprefix(prefix string) {
	f.mutex.Lock()
	f.prefix = prefix
	f.mutex.Unlock()
}
