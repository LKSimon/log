package log

//设置日志flag，即标准库日志flag
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

//设置文件最大大小,主要用于按文件大小切割文件
func (f *FileLogger) SetFileSize(size int64, unit UNIT) {
	f.mutex.Lock()
	f.maxFileSize = int64(unit) * size
	f.mutex.Unlock()
}
