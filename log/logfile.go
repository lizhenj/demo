package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	DATETIME_FORMAT = "2006-01-02 15:04:05"
	DATE_FORMAT     = "2006-01-02"
)

type loggerFile struct {
	fileName string
	suffix   string
	buffs    []string
	lock     sync.Mutex
	flock    sync.Mutex
}

func (l *loggerFile) pop() (datas []string) {
	if l == nil {
		return
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	count := len(l.buffs)
	if count == 0 {
		return
	}
	datas = l.buffs[:]
	l.buffs = l.buffs[count:]

	return
}

func (l *loggerFile) flush() {
	if l == nil {
		return
	}

	datas := l.pop()
	if len(datas) == 0 {
		return
	}

	l.flock.Lock()
	defer l.flock.Unlock()

	var (
		file   *os.File
		writer *bufio.Writer
		err    error
		last   string
	)
	for _, data := range datas {
		tick := data[0:len(DATE_FORMAT)]
		if tick != last {
			if file != nil {
				writer.Flush()
				file.Close()
			}
			file, err = os.OpenFile(fmt.Sprintf("%s_%s.%s", l.fileName, tick, l.suffix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			writer = bufio.NewWriterSize(file, defaultBufSize)
			last = tick
		}
		writer.WriteString(data[:])
	}
	writer.Flush()
	file.Close()
}

func (l *loggerFile) write(buffer string) {
	if l == nil {
		return
	}

	l.lock.Lock()
	l.buffs = append(l.buffs, buffer)
	l.lock.Unlock()
}
