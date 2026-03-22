package logger

import "log"

// SimpleLogger реализует интерфейс Logger для централизованного логирования
// В production заменить на structured logger (zap, zerolog, logrus)
type SimpleLogger struct{}

func (l *SimpleLogger) Info(args ...interface{}) {
	log.Println("INFO:", args)
}

func (l *SimpleLogger) Error(args ...interface{}) {
	log.Println("ERROR:", args)
}
