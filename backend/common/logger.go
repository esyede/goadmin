package common

import (
	"github.com/esyede/goadmin/backend/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.SugaredLogger

func InitLogger() {
	now := time.Now()
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	var coreArr []zapcore.Core

	// Get encoder
	// encoderConfig := zap.NewProductionEncoderConfig()
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Specify time format
	// encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // If not needed, just use zapcore.CapitalLevelEncoder.
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder       // Show full file path
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		// EncodeTime: zapcore.ISO8601TimeEncoder, // ISO8601 UTC time format
		// EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendInt64(int64(d) / 1000000)
		// },
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		// EncodeCaller: zapcore.FullCallerEncoder,
		// EncodeName:       nil,
		// ConsoleSeparator: "",
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// log level
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	// When the level in the yml configuration is greater than Error, the lowPriority level log stops recording.
	if config.Conf.Logs.Level >= 2 {
		lowPriority = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return false
		})
	}

	// info file writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,             // The directory where log files are stored. If the folder does not exist, it will be created automatically.
		MaxSize:    config.Conf.Logs.MaxSize,    // File size limit, unit MB
		MaxAge:     config.Conf.Logs.MaxAge,     // Number of days to retain log files
		MaxBackups: config.Conf.Logs.MaxBackups, // Maximum number of log files to keep
		LocalTime:  false,                       // Do not use local time
		Compress:   config.Conf.Logs.Compress,   // Whether to compress or not
	})
	// The third and subsequent parameters are the log levels written to the file. ErrorLevel mode only records error level logs.
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	// Error file writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,            // Log file storage directory
		MaxSize:    config.Conf.Logs.MaxSize,    // File size limit, unit MB
		MaxAge:     config.Conf.Logs.MaxAge,     // Number of days to retain log files
		MaxBackups: config.Conf.Logs.MaxBackups, // Maximum number of log files to keep
		LocalTime:  false,                       // Do not use local time
		Compress:   config.Conf.Logs.Compress,   // Whether to compress or not
	})
	// The third and subsequent parameters are the log levels written to the file. ErrorLevel mode only records error level logs.
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	Log = logger.Sugar()
	Log.Info("Initialization of zap log completed!")
}
