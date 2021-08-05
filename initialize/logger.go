package initialize

import (
	"crontab/global"
	"crontab/utils"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	// 注意global.Settings.LogsAddress是在settings-dev.yaml配置过的
	// 配置日志的输出地址
	cfg.OutputPaths = []string{
		fmt.Sprintf("%slog_%s.log", global.Settings.LogsAddress, utils.GetNowFormatTodayTime()), //
		"stdout",
	}
	//配置编码方式
	cfg.Encoding = "json"
	cfg.EncoderConfig = encodeConfig()
	// 创建logger实例
	logg, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	global.Lg = logg         // 注册到全局变量中
}

func encodeConfig() zapcore.EncoderConfig {
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(global.TM_FMT_WITH_MS) + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + global.TraceId + "]")
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}
	return zapcore.EncoderConfig{
		CallerKey:      "caller_line", // 打印文件名和行数
		LevelKey:       "level_name",
		MessageKey:     "msg",
		TimeKey:        "ts",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,   // 自定义时间格式
		EncodeLevel:    customLevelEncoder,  // 小写编码器
		EncodeCaller:   customCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}
