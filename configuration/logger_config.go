package configuration

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// Init to example by default in case InitMethod is not called
var Logger *zap.SugaredLogger = zap.NewExample().Sugar()

func InitLogger(logFilePath string) {
	rawJSON := []byte(`{
		"level": "info",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase",
		  "timeKey": "time",
		  "timeEncoder": "ISO8601"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.OutputPaths = append(cfg.OutputPaths, logFilePath)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger.Sugar()
	defer Logger.Sync()

	Logger.Info("logger construction succeeded")
	Logger.Error(time.Now().Format(time.RFC822Z))
}
