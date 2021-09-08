package xlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

//日志级别
const (
	XLOG_FLAG_DEBUG = 1
	XLOG_FLAG_TRACE = 2
	XLOG_FLAG_INFO  = 4
	XLOG_FLAG_ERROR = 8
	XLOG_FLAG_ALL   = XLOG_FLAG_DEBUG | XLOG_FLAG_TRACE | XLOG_FLAG_INFO | XLOG_FLAG_ERROR
)

//单文件最大日志条数
const XLOG_MAX_LOG_COUNT = 3200000

type XLog struct {
	XLogMaxCount int //日志单个文件最大条数

	mutex   *sync.Mutex
	logName string
	logFile *os.File
	now     *time.Time
	name    string
	path    string
	flag    int
	count   int    //当前日志的条数
	round   int    //当前日志文件数
	errors  uint32 //记录的错误日志数量
}

func NewXLog(path, name string, flag int) *XLog {
	log := new(XLog)

	log.flag = flag
	log.mutex = new(sync.Mutex)
	log.name = name
	log.path = path
	log.XLogMaxCount = XLOG_MAX_LOG_COUNT

	return log
}

func backtrace(skip int) string {
	msg := "Backtrace (most recent call last):\n"
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i + skip)
		if !ok {
			break
		}

		msg = fmt.Sprintf("%s  %3d. %s() %s:%d\n", msg, i, runtime.FuncForPC(pc).Name(), file, line)
	}

	return msg
}

func (m *XLog) GetLogname() string {
	return m.logName
}

func (m *XLog) GetLogErrors() uint32 {
	return m.errors
}

func (m *XLog) prepare() {
	now := time.Now()

	if m.now != nil {
		y1, m1, d1 := m.now.Date()
		y2, m2, d2 := now.Date()

		if y1 == y2 && m1 == m2 && d1 == d2 {
			m.now = &now
			m.count++

			if m.round == m.count/m.XLogMaxCount {
				return
			}
		} else {
			m.count = 0
		}

		m.logFile.Close()
	}

	m.now = &now
	m.round = m.count / m.XLogMaxCount

	d, _ := os.Stat(m.path)
	if d == nil {
		os.Mkdir(m.path, 0755)
	}

	m.logName = fmt.Sprintf("%s-%s-%d.log", path.Join(m.path, m.name), now.Format("20060102"), m.round)
	m.logFile, _ = os.OpenFile(m.logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	fmt.Println(m.logName)
}

func (m *XLog) Close() {
	m.mutex.Lock()
	if m.logFile != nil {
		m.logFile.Close()
	}
	m.mutex.Unlock()
}

func (m *XLog) write(lv, msg string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.prepare()

	Y, M, D := m.now.Date()
	h, mon, s := m.now.Clock()
	content := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d] [%s] %s\n", Y, M, D, h, mon, s, lv, msg)

	if m.logFile != nil {
		m.logFile.WriteString(content)
	}
}

func (m *XLog) Debug(str string, a ...interface{}) {
	if m.flag&XLOG_FLAG_DEBUG == 0 {
		return
	}

	m.write("DEBUG", fmt.Sprintf(str, a...))
}

func (m *XLog) Error(str string, a ...interface{}) {
	if m.flag&XLOG_FLAG_ERROR == 0 {
		return
	}

	m.errors++

	m.write("ERROR", fmt.Sprintf(str, a...))

	if m.flag&XLOG_FLAG_DEBUG != 0 {
		m.Debug(backtrace(0))
	}
}

func (m *XLog) Info(str string, a ...interface{}) {
	if m.flag&XLOG_FLAG_INFO == 0 {
		return
	}

	m.write("INFO", fmt.Sprintf(str, a...))
}

func (m *XLog) Trace(str string, a ...interface{}) {
	if m.flag&XLOG_FLAG_TRACE == 0 {
		return
	}

	m.write("TRACE", fmt.Sprintf(str, a...))
}
