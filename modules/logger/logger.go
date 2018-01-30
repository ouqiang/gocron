package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cihub/seelog"
	"gopkg.in/macaron.v1"
)

// 日志库

type Level int8

var logger seelog.LoggerInterface

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

func InitLogger() {
	config := getLogConfig()
	l, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		panic(err)
	}
	logger = l
}

func Debug(v ...interface{}) {
	if macaron.Env != macaron.DEV {
		return
	}
	write(DEBUG, v)
}

func Debugf(format string, v ...interface{}) {
	if macaron.Env != macaron.DEV {
		return
	}
	writef(DEBUG, format, v...)
}

func Info(v ...interface{}) {
	write(INFO, v)
}

func Infof(format string, v ...interface{}) {
	writef(INFO, format, v...)
}

func Warn(v ...interface{}) {
	write(WARN, v)
}

func Warnf(format string, v ...interface{}) {
	writef(WARN, format, v...)
}

func Error(v ...interface{}) {
	write(ERROR, v)
}

func Errorf(format string, v ...interface{}) {
	writef(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	write(FATAL, v)
}

func Fatalf(format string, v ...interface{}) {
	writef(FATAL, format, v...)
}

func write(level Level, v ...interface{}) {
	defer logger.Flush()

	content := ""
	if macaron.Env == macaron.DEV {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			content = fmt.Sprintf("#%s#%s#%d行#", file, runtime.FuncForPC(pc).Name(), line)
		}
	}

	switch level {
	case DEBUG:
		logger.Debug(content, v)
	case INFO:
		logger.Info(content, v)
	case WARN:
		logger.Warn(content, v)
	case FATAL:
		logger.Critical(content, v)
		os.Exit(1)
	case ERROR:
		logger.Error(content, v)
	}
}

func writef(level Level, format string, v ...interface{}) {
	defer logger.Flush()

	content := ""
	if macaron.Env == macaron.DEV {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			content = fmt.Sprintf("#%s#%s#%d行#", file, runtime.FuncForPC(pc).Name(), line)
		}
	}

	format = content + format

	switch level {
	case DEBUG:
		logger.Debugf(format, v...)
	case INFO:
		logger.Infof(format, v...)
	case WARN:
		logger.Warnf(format, v...)
	case FATAL:
		logger.Criticalf(format, v...)
		os.Exit(1)
	case ERROR:
		logger.Errorf(format, v...)
	}
}

func getLogConfig() string {
	config := `
    <seelog>
        <outputs formatid="main">
            %s
            <filter levels="info,critical,error,warn">
                <file path="log/cron.log" />
            </filter>
        </outputs>
        <formats>
            <format id="main" format="%%Date/%%Time [%%LEV] %%Msg%%n"/>
        </formats>
    </seelog>`

	consoleConfig := ""
	if macaron.Env == macaron.DEV {
		consoleConfig =
			`
            <filter levels="info,debug,critical,warn,error">
                <console />
            </filter>
         `
	}
	config = fmt.Sprintf(config, consoleConfig)

	return config
}
