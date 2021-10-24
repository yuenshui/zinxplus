package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SessionId() string {
	nanoMd5 := fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))))
	return strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(nanoMd5+strconv.FormatInt(time.Now().UnixNano(), 10)))))
}
