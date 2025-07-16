package pip

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a custom logger for pip operations
type Logger struct {
	level      LogLevel
	output     io.Writer
	prefix     string
	flags      int
	enableFile bool
	logFile    *os.File
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level      LogLevel
	Output     io.Writer
	Prefix     string
	EnableFile bool
	LogFile    string
	Flags      int
}

// NewLogger creates a new logger instance
func NewLogger(config *LoggerConfig) (*Logger, error) {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	logger := &Logger{
		level:      config.Level,
		output:     config.Output,
		prefix:     config.Prefix,
		flags:      config.Flags,
		enableFile: config.EnableFile,
	}

	// Set up file logging if enabled
	if config.EnableFile && config.LogFile != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.logFile = file

		// Use multi-writer to write to both stdout and file
		if config.Output != nil {
			logger.output = io.MultiWriter(config.Output, file)
		} else {
			logger.output = file
		}
	}

	if logger.output == nil {
		logger.output = os.Stdout
	}

	return logger, nil
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      LogLevelInfo,
		Output:     os.Stdout,
		Prefix:     "[pip-sdk]",
		EnableFile: false,
		Flags:      log.LstdFlags,
	}
}

// Close closes the logger and any open files
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// IsEnabled checks if a log level is enabled
func (l *Logger) IsEnabled(level LogLevel) bool {
	return level >= l.level
}

// log writes a log message with the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	var caller string
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	} else {
		caller = "unknown"
	}

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Format message
	message := fmt.Sprintf(format, args...)

	// Create log entry
	logEntry := fmt.Sprintf("%s [%s] %s %s - %s\n",
		timestamp,
		level.String(),
		l.prefix,
		caller,
		message,
	)

	// Write to output
	if l.output != nil {
		l.output.Write([]byte(logEntry))
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, format, args...)
	os.Exit(1)
}

// LogCommand logs a command execution
func (l *Logger) LogCommand(command string, args []string, duration time.Duration) {
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}

	l.Info("Command executed: %s (duration: %v)", fullCommand, duration)
}

// LogCommandError logs a command execution error
func (l *Logger) LogCommandError(command string, args []string, err error, output string) {
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}

	l.Error("Command failed: %s", fullCommand)
	l.Error("Error: %v", err)
	if output != "" {
		l.Error("Output: %s", output)
	}
}

// LogPackageOperation logs a package operation
func (l *Logger) LogPackageOperation(operation, packageName string, success bool, duration time.Duration) {
	if success {
		l.Info("Package %s successful: %s (duration: %v)", operation, packageName, duration)
	} else {
		l.Error("Package %s failed: %s (duration: %v)", operation, packageName, duration)
	}
}

// LogVenvOperation logs a virtual environment operation
func (l *Logger) LogVenvOperation(operation, path string, success bool) {
	if success {
		l.Info("Virtual environment %s successful: %s", operation, path)
	} else {
		l.Error("Virtual environment %s failed: %s", operation, path)
	}
}

// LogProjectOperation logs a project operation
func (l *Logger) LogProjectOperation(operation, path string, success bool) {
	if success {
		l.Info("Project %s successful: %s", operation, path)
	} else {
		l.Error("Project %s failed: %s", operation, path)
	}
}

// ParseLogLevel parses a log level string
func ParseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LogLevelDebug, nil
	case "INFO":
		return LogLevelInfo, nil
	case "WARN", "WARNING":
		return LogLevelWarn, nil
	case "ERROR":
		return LogLevelError, nil
	case "FATAL":
		return LogLevelFatal, nil
	default:
		return LogLevelInfo, fmt.Errorf("invalid log level: %s", level)
	}
}

// Update Manager to use the new logger
func (m *Manager) SetCustomLogger(logger *Logger) {
	m.customLogger = logger
}

// logInfo logs an info message using custom logger if available
func (m *Manager) logInfo(format string, args ...interface{}) {
	if m.customLogger != nil {
		m.customLogger.Info(format, args...)
	} else if m.logger != nil {
		m.logger.Printf("[INFO] "+format, args...)
	}
}

// logError logs an error message using custom logger if available
func (m *Manager) logError(format string, args ...interface{}) {
	if m.customLogger != nil {
		m.customLogger.Error(format, args...)
	} else if m.logger != nil {
		m.logger.Printf("[ERROR] "+format, args...)
	}
}

// logDebug logs a debug message using custom logger if available
func (m *Manager) logDebug(format string, args ...interface{}) {
	if m.customLogger != nil {
		m.customLogger.Debug(format, args...)
	} else if m.logger != nil && m.config.LogLevel == "DEBUG" {
		m.logger.Printf("[DEBUG] "+format, args...)
	}
}

// logWarn logs a warning message using custom logger if available
func (m *Manager) logWarn(format string, args ...interface{}) {
	if m.customLogger != nil {
		m.customLogger.Warn(format, args...)
	} else if m.logger != nil {
		m.logger.Printf("[WARN] "+format, args...)
	}
}
