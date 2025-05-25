package logger

import (
	"fmt"
	"godoku/config"
	"log"
	"os"
)

var (
	// Logger is the global logger instance
	Logger *log.Logger
	// LogFile is the current log file
	LogFile *os.File
)

// InitLogger initializes the logger using the log file path from config
func InitLogger() error {
	// Close any existing log file
	if LogFile != nil {
		LogFile.Close()
	}

	// Ensure log directory exists
	if err := config.InitLogDirectory(); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Get log file path from config
	logFilePath := config.GetLogPath()
	
	var err error
	LogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Set the global logger
	Logger = log.New(LogFile, "", log.LstdFlags|log.Lshortfile)
	Logger.Printf("Logging initialized at %s", logFilePath)
	
	return nil
}

// Close closes the log file
func Close() {
	if LogFile != nil {
		LogFile.Close()
	}
}

// Info logs an informational message
func Info(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[DEBUG] "+format, v...)
	}
}
