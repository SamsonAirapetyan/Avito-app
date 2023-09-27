package logger

import (
	"github.com/hashicorp/go-hclog"
	"os"
	"sync"
)

var logger hclog.Logger

// sync.Once служит для того, чтобы функция выполнилась лишь раз, несмотря на то, сколько раз она будет вызвана
// Материал для просмотра https://golang-blog.blogspot.com/2019/10/sync-once-do-go.html
var once sync.Once

// GetLogger return new Logger
func GetLogger() hclog.Logger {
	once.Do(func() {
		loggerOption := &hclog.LoggerOptions{
			Name:            "Concurrency",
			Level:           hclog.Info,
			Output:          os.Stderr,
			IncludeLocation: true,
			TimeFormat:      "2006-01-02 15:04:05.000",
		}
		logger = hclog.New(loggerOption)
	})
	return logger
}
