package initialization

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	once     sync.Once
	instance *zap.Logger
)

// InitLogger 初始化全局日志记录器
func InitLogger(opts ...Option) *zap.Logger {
	once.Do(func() {
		// 默认配置
		config := Config{
			LogLevel:      zapcore.DebugLevel,
			LogFile:       "./logs/app.log",
			ConsoleOutput: true,
			FileOutput:    true,
			MaxSize:       100, // MB
			MaxBackups:    5,
			MaxAge:        30, // days
			Buffered:      true,
			BufferSize:    4096,
		}

		// 应用选项模式配置
		for _, opt := range opts {
			opt(&config)
		}

		// 创建核心
		core := zapcore.NewCore(
			getEncoder(),
			getWriteSyncer(config),
			config.LogLevel,
		)

		// 创建 logger 并添加 caller 信息
		instance = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)

		// 替换全局 logger
		zap.ReplaceGlobals(instance)
	})

	return instance
}

// Config 日志配置结构
type Config struct {
	LogLevel      zapcore.Level
	LogFile       string
	ConsoleOutput bool
	FileOutput    bool
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	Buffered      bool
	BufferSize    int
}

// Option 配置选项函数类型
type Option func(*Config)

// WithLogLevel 设置日志级别
func WithLogLevel(level zapcore.Level) Option {
	return func(c *Config) {
		c.LogLevel = level
	}
}

// WithLogFile 设置日志文件路径
func WithLogFile(path string) Option {
	return func(c *Config) {
		c.LogFile = path
	}
}

// WithFileRotation 设置日志轮转配置
func WithFileRotation(maxSize, maxBackups, maxAge int, compress bool) Option {
	return func(c *Config) {
		c.MaxSize = maxSize
		c.MaxBackups = maxBackups
		c.MaxAge = maxAge
		c.Compress = compress
	}
}

// WithOutput 设置输出目标
func WithOutput(console, file bool) Option {
	return func(c *Config) {
		c.ConsoleOutput = console
		c.FileOutput = file
	}
}

// WithBuffered 设置缓冲配置
func WithBuffered(enabled bool, size int) Option {
	return func(c *Config) {
		c.Buffered = enabled
		c.BufferSize = size
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer(config Config) zapcore.WriteSyncer {
	var syncers []zapcore.WriteSyncer

	if config.ConsoleOutput {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if config.FileOutput {
		// 创建日志目录
		if err := os.MkdirAll("./logs", 0755); err != nil {
			zap.L().Error("Failed to create log directory", zap.Error(err))
		}

		// 使用 lumberjack 实现日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		ws := zapcore.AddSync(fileWriter)

		// 添加缓冲层
		if config.Buffered {
			ws = &zapcore.BufferedWriteSyncer{
				WS:            ws,
				Size:          config.BufferSize,
				FlushInterval: 5 * time.Second,
			}
		}

		syncers = append(syncers, ws)
	}

	if len(syncers) == 0 {
		// 默认至少输出到控制台
		return zapcore.AddSync(os.Stdout)
	}

	return zapcore.NewMultiWriteSyncer(syncers...)
}
