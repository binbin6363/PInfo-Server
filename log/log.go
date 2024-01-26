package log

import (
	"context"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

var defaultLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	return defaultLogger
}

func InitLogger(fileName string, maxSize, maxBackups, maxAge, level, callerSkip int) {
	writeSyncer := getLogWriter(fileName, maxSize, maxBackups, maxAge)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(level))

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
	defaultLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileName string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize, // mb
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // day
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

type ZapLogger struct {
	logger *zap.Logger
}

func (l ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l ZapLogger) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Info(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

func (l ZapLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Warn(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

func (l ZapLogger) Error(ctx context.Context, s string, i ...interface{}) {
	l.logger.Error(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

func (l ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	fields := getContextFields(ctx)
	if err != nil {
		sql, rows := fc()
		l.logger.Error(err.Error(), append(fields, zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))...)
	} else {
		sql, rows := fc()
		l.logger.Debug("", append(fields, zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))...)
	}
}

func getContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	// 将 ctx 中的字段添加到 fields 切片中，根据需要自定义
	// 例如，你可以将请求 ID、用户 ID 等上下文信息添加到日志中
	// fields = append(fields, zap.String("requestID", getRequestIDFromContext(ctx)))
	return fields
}

func NewZapLogger() logger.Interface {
	return ZapLogger{logger: GetLogger().Desugar()}
}
