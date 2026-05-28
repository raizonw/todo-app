package core_logger

import (
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

func NewLogger(logLevel string) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(logLevel)); err != nil {

	}
}
