package mysqlx

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/IguoChan/go-pkg/logger"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Logger struct {
	*glogger.Config

	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(config *glogger.Config) *Logger {
	var (
		traceStr     = "%s: [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s: [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s: [%.3fms] [rows:%v] %s"
	)

	return &Logger{
		Config:       config,
		traceStr:     traceStr,
		traceErrStr:  traceWarnStr,
		traceWarnStr: traceErrStr,
	}
}

func (l *Logger) LogMode(level glogger.LogLevel) glogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= glogger.Info {
		logger.Infof(ctx, s, i...)
	}
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= glogger.Warn {
		logger.Infof(ctx, s, i...)
	}
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= glogger.Error {
		logger.Infof(ctx, s, i...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= glogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= glogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.Errorf(ctx, l.traceErrStr, FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Errorf(ctx, l.traceErrStr, FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= glogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Warnf(ctx, l.traceWarnStr, FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Warnf(ctx, l.traceWarnStr, FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == glogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Infof(ctx, l.traceStr, FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Infof(ctx, l.traceStr, FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

var (
	sourceDir string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get mysqlx source directory with various operating systems
	sourceDir = regexp.MustCompile(`mysqlx.logger\.go`).ReplaceAllString(file, "")
}

// FileWithLineNum
// @description: 参照 utils.FileWithLineNum 进行改造
// @return string
func FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		// 去除dbx中的调用栈
		// 去除 gorm 文件中的调用栈
		if ok && (!strings.HasPrefix(file, sourceDir) && !strings.Contains(file, "gorm.io/gorm") || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
