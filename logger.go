package svrkit

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LOG_EMERG = 0
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARN
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG = 7
)

type Logger_t struct {
	priority int
	l        *log.Logger
	m        sync.RWMutex
}

var logger *Logger_t

func NewLogger(prefix string, priority int) {
	logger = new(Logger_t)
	logger.priority = priority

	var logFile *os.File

	now := time.Now()

	logger.m.Lock()
	logFileName := fmt.Sprintf("%v_%04v%02v%02v.log", prefix, now.Year(), int(now.Month()), now.Day())
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	logger.l = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.m.Unlock()

	go func() {
		for {
			nextTick := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			if !nextTick.After(now) {
				nextTick = nextTick.Add(24 * time.Hour)
			}
			diff := nextTick.Sub(now)
			ticker := time.NewTicker(diff)
			<-ticker.C

			now = time.Now()

			logger.m.Lock()
			if logFile != nil {
				err = logFile.Close()
				if err != nil {
					panic(err)
				}
				logFile = nil
			}
			logFileName := fmt.Sprintf("%v_%04v%02v%02v.log", prefix, now.Year(), int(now.Month()), now.Day())
			logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
			if err != nil {
				panic(err)
			}
			logger.l = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
			logger.m.Unlock()

		}
	}()

	return
}

func SetPriority(priority int) {
	logger.m.Lock()
	logger.priority = priority
	logger.m.Unlock()
}

func GetPriority() int {
	logger.m.RLock()
	priority := logger.priority
	logger.m.RUnlock()
	return priority
}

func Emerg(v ...interface{}) {
	if logger.priority >= LOG_EMERG {
		logger.m.RLock()
		logger.l.Print(v...)
		logger.m.RUnlock()
	}
}

func Emergf(format string, v ...interface{}) {
	if logger.priority >= LOG_EMERG {
		logger.m.RLock()
		logger.l.Printf(format, v...)
		logger.m.RUnlock()
	}
}

func Alert(v ...interface{}) {
	if logger.priority >= LOG_ALERT {
		logger.m.RLock()
		logger.l.Print(v...)
		logger.m.RUnlock()
	}
}

func Alertf(format string, v ...interface{}) {
	if logger.priority >= LOG_ALERT {
		logger.m.RLock()
		logger.l.Printf(format, v...)
		logger.m.RUnlock()
	}
}

func Crit(v ...interface{}) {
	if logger.priority >= LOG_CRIT {
		logger.m.RLock()
		logger.l.Print(v...)
		logger.m.RUnlock()
	}
}

func Critf(format string, v ...interface{}) {
	if logger.priority >= LOG_CRIT {
		logger.m.RLock()
		logger.l.Printf(format, v...)
		logger.m.RUnlock()
	}
}

func Err(v ...interface{}) {
	if logger.priority >= LOG_ERR {
		logger.m.RLock()
		logger.l.Print(append([]interface{}{"[ERR] "}, v...)...)
		logger.m.RUnlock()
	}
}

func Errf(format string, v ...interface{}) {
	if logger.priority >= LOG_ERR {
		logger.m.RLock()
		logger.l.Printf("[ERR] "+format, v...)
		logger.m.RUnlock()
	}
}

func Warn(v ...interface{}) {
	if logger.priority >= LOG_WARN {
		logger.m.RLock()
		logger.l.Print(append([]interface{}{"[WARN] "}, v...)...)
		logger.m.RUnlock()
	}
}

func Warnf(format string, v ...interface{}) {
	if logger.priority >= LOG_WARN {
		logger.m.RLock()
		logger.l.Printf("[WARN] "+format, v...)
		logger.m.RUnlock()
	}
}

func Notice(v ...interface{}) {
	if logger.priority >= LOG_NOTICE {
		logger.m.RLock()
		logger.l.Print(append([]interface{}{"[NOTICE] "}, v...)...)
		logger.m.RUnlock()
	}
}

func Noticef(format string, v ...interface{}) {
	if logger.priority >= LOG_NOTICE {
		logger.m.RLock()
		logger.l.Printf("[NOTICE] "+format, v...)
		logger.m.RUnlock()
	}
}

func Info(v ...interface{}) {
	if logger.priority >= LOG_INFO {
		logger.m.RLock()
		logger.l.Print(append([]interface{}{"[INFO] "}, v...)...)
		logger.m.RUnlock()
	}
}

func Infof(format string, v ...interface{}) {
	if logger.priority >= LOG_INFO {
		logger.m.RLock()
		logger.l.Printf("[INFO] "+format, v...)
		logger.m.RUnlock()
	}
}

func Debug(v ...interface{}) {
	if logger.priority >= LOG_DEBUG {
		logger.m.RLock()
		logger.l.Print(append([]interface{}{"[DEBUG] "}, v...)...)
		logger.m.RUnlock()
	}
}

func Debugf(format string, v ...interface{}) {
	if logger.priority >= LOG_DEBUG {
		logger.m.RLock()
		logger.l.Printf("[DEBUG] "+format, v...)
		logger.m.RUnlock()
	}
}
