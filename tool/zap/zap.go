package logger

import (
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

var logger *zap.Logger

// InitLog 初始化日志 logger
func InitLog(logPath, errPath string, level string, locText func(MessageIDs ...string) string) {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, //将级别转换成大写
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	encoder := zapcore.NewConsoleEncoder(config)
	// 设置级别
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	// 实现两个判断日志等级的interface  可以自定义级别展示
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现

	var zapCores []zapcore.Core
	var infoWriter, warnWriter io.Writer
	var err error
	// 将info及以下写入logPath,  warn及以上写入errPath
	if logPath != "" {
		infoWriter, err = getWriter(logPath)
		zapCores = append(zapCores, zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel))
	}
	if errPath != "" {
		warnWriter, err = getWriter(errPath)
		zapCores = append(zapCores, zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel))
	}
	if err != nil {
		log.Println(locText("loggingSystemStartupException"))
		panic(err)
	}
	//日志都会在console中展示
	zapCores = append(zapCores, zapcore.NewCore(zapcore.NewConsoleEncoder(config),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel))

	// 最后创建具体的Logger
	core := zapcore.NewTee(zapCores...)
	logger = zap.New(core)
	//logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
}

func getWriter(filename string) (io.Writer, error) {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	hook, err := rotateLogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotateLogs.WithLinkName(filename),
		rotateLogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotateLogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)

	return hook, err
}

// logs.Debug(...)
func Debug(format string, v ...interface{}) {
	logger.Sugar().Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Sugar().Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Sugar().Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Sugar().Errorf(format, v...)
}

func Panic(format string, v ...interface{}) {
	logger.Sugar().Panicf(format, v...)
}

func DropErr(err error) {
	if err != nil {
		Panic("%w", err)
	}
}
