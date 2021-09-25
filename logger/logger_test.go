package logger_test

import (
	"github/driftingboy/tsdb/logger"
	"os"
	"testing"
	"time"
)

func Test_Logger(t *testing.T) {
	log := logger.NewStdLogger(os.Stdout, 4096)
	log.Log(logger.Error, "instance", "119.3.63.51", "empty")
}

func Test_WithPrefix(t *testing.T) {
	std := logger.NewStdLogger(os.Stdout, 4096)
	log := logger.WithPrefix(std, logger.DefaultCaller, logger.Timer(time.RFC3339))
	log.Log(logger.Debug, "instance", "119.3.63.51")
}
