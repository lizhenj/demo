package base

import (
	"fmt"
	"runtime"
)

func FileLine(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	i, count := len(file)-4, 0
	for ; i > 0; i-- {
		if file[i] == '/' {
			count++
			if count == 2 {
				break
			}
		}
	}
	return fmt.Sprintf("%s:%d", file[i+1:], line)
}
