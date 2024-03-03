package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	operationPath string
	errorPath     string
)

var operateLog *zap.SugaredLogger
var errorLog *zap.SugaredLogger

func GetOperateLog() *zap.SugaredLogger {
	return operateLog
}

func GetErrorLog() *zap.SugaredLogger {
	return errorLog
}

// SetupZapLogger 启动Zap记录器
func SetupZapLogger() {
	operationPath = Conf.LogPath.OperationPath
	errorPath = Conf.LogPath.ErrorPath
	encoder := getEncoder()
	writeSyncer1 := getWriteSyncer(operationPath)
	operateLog = zap.New(zapcore.NewCore(encoder, writeSyncer1, zapcore.DebugLevel)).Sugar()
	writeSyncer2 := getWriteSyncer(errorPath)
	errorLog = zap.New(zapcore.NewCore(encoder, writeSyncer2, zapcore.DebugLevel)).Sugar()
}

func getEncoder() zapcore.Encoder {
	// 覆盖NewProductionEncoderConfig()默认的编码设置，改为自定义的
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 以json格式的方式输出日志
	//encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// 以普通方式输出日志
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return encoder
}

func getWriteSyncer(filepath string) zapcore.WriteSyncer {
	if !IsExist(filepath) {
		CreatFile(filepath)
	}

	logger := lumberjack.Logger{
		Filename: filepath,
		// 日志文件每1MB会切割并且在当前目录下最多保存5个备份
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		//参数含义
		//Filename: 日志文件的位置
		//MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
		//MaxBackups：保留旧文件的最大个数
		//MaxAges：保留旧文件的最大天数
		//Compress：是否压缩/归档旧文件
	}
	return zapcore.AddSync(&logger)
}
