package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

type LevelWriter struct {
	io.Writer
	Levels []zerolog.Level
}

func init() {
	stdoutWriter := LevelWriter{Writer: os.Stdout, Levels: []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel}}
	stderrWriter := LevelWriter{Writer: os.Stderr, Levels: []zerolog.Level{zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel}}

	multi := zerolog.MultiLevelWriter(
		stdoutWriter,
		stderrWriter,
	)

	log = zerolog.New(multi)
}

func Instance() *zerolog.Logger {
	return &log
}

func (lw LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	for _, l := range lw.Levels {
		if l == level {
			return lw.Write(p)
		}
	}
	return len(p), nil
}
