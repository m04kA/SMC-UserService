package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	info    *log.Logger
	warn    *log.Logger
	err     *log.Logger
	logFile *os.File
}

var std *Logger

// Init инициализирует логгер с записью в консоль и файл
func Init(logFilePath string) error {
	// Открываем файл для логов
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Настраиваем мультиплексор для warning и error
	warnErrorWriter := io.MultiWriter(os.Stdout, logFile)

	std = &Logger{
		info:    log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		warn:    log.New(warnErrorWriter, "[WARN] ", log.LstdFlags|log.Lshortfile),
		err:     log.New(warnErrorWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		logFile: logFile,
	}

	return nil
}

// Close закрывает файл логов
func Close() error {
	if std != nil && std.logFile != nil {
		return std.logFile.Close()
	}
	return nil
}

// Info логирует информационные сообщения (только консоль)
func Info(format string, v ...interface{}) {
	if std != nil {
		std.info.Printf(format, v...)
	}
}

// Warn логирует предупреждения (консоль + файл)
func Warn(format string, v ...interface{}) {
	if std != nil {
		std.warn.Printf(format, v...)
	}
}

// Error логирует ошибки (консоль + файл)
func Error(format string, v ...interface{}) {
	if std != nil {
		std.err.Printf(format, v...)
	}
}

// Fatal логирует критическую ошибку и завершает программу
func Fatal(format string, v ...interface{}) {
	if std != nil {
		std.err.Fatalf(format, v...)
	}
}
