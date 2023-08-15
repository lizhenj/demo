package log

import (
	"demo/base"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

type LogLevel = int

const (
	ERROR_N LogLevel = iota
	INFO_N
	WARN_N
)

const (
	defaultBufSize = 1024 * 10240
	skipLevel      = 3
)

type loggerPrint struct {
	buffs []string
	lock  sync.Mutex
	flock sync.Mutex
}

func (l *loggerPrint) pop() (data []string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	count := len(l.buffs)
	if count == 0 {
		return
	}
	data = l.buffs[:]
	l.buffs = l.buffs[count:]
	return
}

func (l *loggerPrint) flush() {
	datas := l.pop()
	if len(datas) == 0 {
		return
	}

	l.flock.Lock()
	defer l.flock.Unlock()

	for _, data := range datas {
		fmt.Printf(data)
	}
}

func (l *loggerPrint) write(buffer string) {
	l.lock.Lock()
	l.buffs = append(l.buffs, buffer)
	l.lock.Unlock()
}

type logger struct {
	//日志输出
	print *loggerPrint
	//日志存储
	record *loggerFile
	fatal  *loggerFile
}

var (
	Logger = newLogger()

	SetFile = Logger.SetFile
	Close   = Logger.Close

	Info   = Logger.Info
	Infof  = Logger.Infof
	Error  = Logger.Error
	Errorf = Logger.Errorf
)

var (
	levelText = map[int]string{
		ERROR_N: "ERROR",
		INFO_N:  "INFO",
		WARN_N:  "WARN",
	}
)

func newLogger() *logger {
	l := &logger{}
	l.print = &loggerPrint{
		buffs: make([]string, 0),
	}

	go func() {
		var delay time.Duration
		if runtime.GOOS == "windows" {
			delay = time.Millisecond * 100
		} else {
			delay = time.Second
		}
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				l.print.flush()
				t.Reset(delay)
			}
		}
	}()
	return l
}

func (l *logger) SetFile(fileName string) {
	if err := os.MkdirAll(path.Dir(fileName), os.ModePerm); err != nil {
		panic(err)
	}

	l.record = &loggerFile{
		fileName: fileName,
		suffix:   "log",
		buffs:    make([]string, 0),
	}
	l.fatal = &loggerFile{
		fileName: fileName,
		suffix:   "error",
		buffs:    make([]string, 0),
	}

	go func() {
		delay := time.Second
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				l.record.flush()
				l.fatal.flush()
				t.Reset(delay)
			}
		}
	}()
}

func (l *logger) Close() {
	l.print.flush()

	l.record.flush()
	l.fatal.flush()
}

func (l *logger) write(level LogLevel, data string) {
	text := fmt.Sprintf("%s %s [%s] -%s\n", time.Now().Format(DATETIME_FORMAT), levelText[level], base.FileLine(skipLevel), data)

	l.print.write(text)
	l.record.write(text)
	if level == ERROR_N {
		l.fatal.write(text)
	}
}

func (l *logger) Write(p []byte) (n int, err error) {
	text := fmt.Sprintf("%s %s [%s] -%s\n", time.Now().Format(DATETIME_FORMAT), "INFO", base.FileLine(skipLevel), p)
	fmt.Printf("测试测试log 哈哈哈哈 %v", text)
	l.print.write(text)
	l.record.write(text)

	return 0, nil
}

func (l *logger) Error(data ...interface{}) {
	l.write(ERROR_N, fmt.Sprint(data...))
}

func (l *logger) Errorf(format string, data ...interface{}) {
	l.write(ERROR_N, fmt.Sprintf(format, data...))
}

func (l *logger) Info(data ...interface{}) {
	l.write(INFO_N, fmt.Sprint(data...))
}

func (l *logger) Infof(format string, data ...interface{}) {
	l.write(INFO_N, fmt.Sprintf(format, data...))
}

func (l *logger) Print(data ...interface{}) {
	l.write(WARN_N, fmt.Sprint(data...))
}

func (l *logger) Printf(format string, data ...interface{}) {
	l.write(WARN_N, fmt.Sprintf(format, data...))
}
