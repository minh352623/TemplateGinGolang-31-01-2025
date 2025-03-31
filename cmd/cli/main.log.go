package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	endcoder := getEncoderLog()
	sync := getWriteSync()
	core := zapcore.NewCore(endcoder, sync, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger.Info("hello world")
}

// format logs
func getEncoderLog() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}


// write logs to file
func getWriteSync() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.txt",os.O_CREATE|os.O_WRONLY,os.ModePerm)
	syncFile := zapcore.AddSync(file)
	syncConsole := zapcore.AddSync(os.Stderr)
	return zapcore.NewMultiWriteSyncer(syncConsole, syncFile)
}