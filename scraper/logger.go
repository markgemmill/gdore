package scraper

import (
	"fmt"
	"os"
	"time"

	"github.com/markgemmill/pathlib"
)

type Logger struct {
	file *os.File
}

func (l *Logger) Log(msg string) {
	l.file.WriteString(msg)
	l.file.WriteString("\n")
	fmt.Println(msg)
}

func (l *Logger) Logf(msg string, args ...any) {
	smsg := fmt.Sprintf(msg, args...)
	l.Log(smsg)
}

func (l *Logger) Close() {
	l.file.Close()
}

var _logger *Logger

func GetLogger() *Logger {
	return _logger
}

func InitLogger(logDir pathlib.Path) (*Logger, error) {
	timestamp := time.Now().Format("2006-01-02-150405")
	logFile := logDir.Join(fmt.Sprintf("log-%s.txt", timestamp))
	file, err := logFile.Open()
	if err != nil {
		return nil, err
	}
	_logger = &Logger{
		file: file,
	}

	return _logger, nil
}
