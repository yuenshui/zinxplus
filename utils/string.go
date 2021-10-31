package utils

import (
	"bytes"
	"time"
)

//拼接字符串使用方法
//项目中杜绝使用字符串+字符串的形式，杜绝使用[]byte和字符串互相转换拼接的形式
//规范使用，统一使用方法，提高效率，节约资源
//使用方式 1  ConcatString("aaa","bbb","ccc")
//使用方式 2  ConcatString([]string...)
func ConcatString(args ...string) string {

	var buffer bytes.Buffer
	for _, i := range args {
		buffer.WriteString(i)
	}

	return buffer.String()
}

// Microtime 毫秒
func Microtime() int64 {
	return time.Now().UnixNano() / 1e6 // 纳秒转毫秒
}

// Nanotime 纳秒
func Nanotime() int64 {
	return time.Now().UnixNano() // 纳秒
}
