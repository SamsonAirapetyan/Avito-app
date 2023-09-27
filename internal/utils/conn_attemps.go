package utils

import (
	logger2 "SergeyProject/pkg/logger"
	"time"
)

// ConnectionAttemps try to connect to BD attemps time
func ConnectionAttemps(conn_func func() error, attemps int, delay time.Duration) error {
	logger := logger2.GetLogger()
	var err error
	for i := 0; i < attemps; i++ {
		err = conn_func()
		if err != nil {
			logger.Warn("Attempting to connect\"", "current attemp", i+1, "appemps left", attemps-i-1)
			time.Sleep(delay)
			continue
		}
	}
	return err
}
