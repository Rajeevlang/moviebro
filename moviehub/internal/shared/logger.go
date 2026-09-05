package shared

import (
	"go.uber.org/zap"
)

// Log is the global logger instance
var Log *zap.Logger

// InitLogger initializes a production-ready Zap logger.
// Call this early in the main() function of your microservices.
func InitLogger() {
	var err error

	// You can switch to zap.NewDevelopment() for colorful console output during local dev
	Log, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Ensures logs are flushed before the application exits
	defer Log.Sync()
}
