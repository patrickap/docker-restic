package log

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	stdoutWriter := &LevelWriter{Writer: os.Stdout, Levels: []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel}}
	stderrWriter := &LevelWriter{Writer: os.Stderr, Levels: []zerolog.Level{zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel}}

	multi := zerolog.MultiLevelWriter(
		stdoutWriter,
		stderrWriter,
	)

	log = zerolog.New(multi)
}

func Instance() *zerolog.Logger {
	return &log
}

type LevelWriter struct {
	io.Writer
	Levels []zerolog.Level
}

func (lw *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	for _, l := range lw.Levels {
		if l == level {
			return lw.Write(p)
		}
	}
	return len(p), nil
}

type LogWrapper struct {
	Writer io.Writer
	Logger func() *zerolog.Event
}

func (lw *LogWrapper) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		lw.Logger().Msg(scanner.Text())
	}
	return len(p), scanner.Err()
}
