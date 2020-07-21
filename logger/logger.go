package logger

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/logger"
	"github.com/jackc/pgx/v4"
)

type Logger struct {
	lvl    pgx.LogLevel
	prefix string
	*logger.Logger
}

func NewLogger(name, prefix, logPath string, lvl pgx.LogLevel) *Logger {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	return &Logger{
		prefix: prefix,
		Logger: logger.Init(name, true, true, lf),
	}
}

func (lg *Logger) SubLog(prefix string) *Logger {
	return &Logger{
		prefix: lg.prefix + "." + prefix,
		Logger: lg.Logger,
	}
}

func (lg *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	bt, _ := json.Marshal(data)
	if lg.lvl < level {
		return
	}
	switch level {
	case pgx.LogLevelTrace:
		lg.Infof("%s Trace: %s - %s\n", lg.prefix, msg, string(bt))
	case pgx.LogLevelDebug:
		lg.Infof("%s Debug: %s - %s\n", lg.prefix, msg, string(bt))
	case pgx.LogLevelInfo:
		lg.Infof("%s Info: %s - %s\n", lg.prefix, msg, string(bt))
	case pgx.LogLevelWarn:
		lg.Warningf("%s Warn: %s - %s\n", lg.prefix, msg, string(bt))
	case pgx.LogLevelError:
		lg.Errorf("%s Error: %s - %s\n", lg.prefix, msg, string(bt))
	case pgx.LogLevelNone:
		lg.Errorf("%s Epic: %s - %s\n", lg.prefix, msg, string(bt))
	default:
		lg.Errorf("%s ????: %s - %s\n", lg.prefix, msg, string(bt))
	}
}
