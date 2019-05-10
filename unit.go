package log

import (
	"fmt"
	"os"
	"path/filepath"
)

//连接文件路径于文件名
func joinFilePath(dir, fileName string) string {
	if "" == fileName {
		fileName = "log"
	}

	return filepath.Join(dir, fileName)
}

//判断文件或者路径是否存在
func isExist(name string) bool {
	_, err := os.Stat(name)
	if nil == err {
		return true
	}

	if os.IsExist(err) {
		return true
	}

	return false
}

//获取文件大小
func fileSize(name string) int64 {
	fileInfo, err := os.Stat(name)
	if nil != err {
		return 0
	}

	fmt.Println("文件大小为： ", fileInfo.Size())

	return fileInfo.Size()
}

//获取文件名
func getFileName(path string) string {
	if "" == path {
		return "."
	}

	fileName := filepath.Base(path)
	name := getStringBeforeSpecificChar(fileName, '.')

	return name
}

//获取string特定字符前的数据
func getStringBeforeSpecificChar(str string, c byte) string {
	rb := []byte{}
	b := []byte(str)
	for _, v := range b {
		if v != c {
			rb = append(rb, v)
		} else {
			break
		}
	}

	return string(rb)
}
