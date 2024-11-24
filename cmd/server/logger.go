package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"project/internal/conf"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
)

// NewLogger creates a logger based on the configuration
//
// Note that the logger should typically write multiple logs simultaneously to different
// outputs (Standard Output [os.Stdout], Log files, or Log collectors like logstash, filebeat, etc.).
func NewLogger(c *conf.Log) log.Logger {
	return kratoszap.NewLogger(
		zap.New(
			zapcore.NewTee(
				zapcore.NewCore(
					zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
						LevelKey:         "level",
						MessageKey:       "msg",
						TimeKey:          "ts",
						StacktraceKey:    "st",
						LineEnding:       zapcore.DefaultLineEnding,
						ConsoleSeparator: " ",
						EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
							encoder.AppendString(
								fmt.Sprintf(
									"\033[36m%s:%d\033[0m",
									caller.File,
									caller.Line,
								),
							)
						},
						// Print out the timestamp with customized formatting
						EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
							encoder.AppendString(t.Format(timeFormat))
						},
						// Print out different levels using different colors
						EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
							switch level {
							case zapcore.DebugLevel:
								encoder.AppendString("\033[37;1mDEBUG\033[0m")
							case zapcore.InfoLevel:
								encoder.AppendString("\033[36;1mINFO\033[0m")
							case zapcore.WarnLevel:
								encoder.AppendString("\033[33;1mWARN\033[0m")
							case zapcore.ErrorLevel:
								encoder.AppendString("\033[31;1mERROR\033[0m")
							case zapcore.DPanicLevel:
								encoder.AppendString("\033[31;1mDPANIC\033[0m")
							case zapcore.PanicLevel:
								encoder.AppendString("\033[31;1mPANIC\033[0m")
							case zapcore.FatalLevel:
								encoder.AppendString("\033[31;1mFATAL\033[0m")
							case zapcore.InvalidLevel:
								encoder.AppendString("\033[31;1mINVALID\033[0m")
							}
							encoder.AppendString(fmt.Sprintf("\033[35m[%s]\033[0m", id))
						},
					}),
					// We need to tell the zap core that we want the log message to be printed to the console
					zapcore.AddSync(os.Stdout),
					// All the log messages whose level is beyond [zapcore.DebugLevel] are printed
					zapcore.Level(c.Level-1),
				),
				// You can add one more output path as well by uncommenting the following lines.
				zapcore.NewCore(
					zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
					zapcore.AddSync(&lumberjack.Logger{
						Filename:   filepath.Join(c.Path, Name+"-"+Version+".log"),
						MaxAge:     int(c.MaxAge),
						MaxBackups: int(c.MaxBackups),
						MaxSize:    int(c.MaxSize),
					}),
					zapcore.Level(c.Level-1),
				),
			),
			// Print out the stack trace when the level is greater than or equal to [zap.WarnLevel]
			zap.AddStacktrace(zap.WarnLevel),
		),
	)
}
