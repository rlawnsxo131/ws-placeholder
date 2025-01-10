package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	_onceDefaultLogger     sync.Once
	singletonDefaultLogger *DefaultLogger
)

func Default() *DefaultLogger {
	return singletonDefaultLogger
}

type DefaultLogger struct {
	*zerolog.Logger
}

func NewDefaultLogger(w io.Writer) *DefaultLogger {
	// ko: +9 hours
	l := zerolog.New(w).With().Caller().Str("time", time.Now().UTC().Format(time.RFC3339Nano)).Logger()
	return &DefaultLogger{&l}
}

func init() {
	_onceDefaultLogger.Do(func() {
		singletonDefaultLogger = NewDefaultLogger(os.Stdout)
	})
}
