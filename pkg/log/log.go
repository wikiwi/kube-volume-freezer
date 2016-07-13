package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	restfullog "github.com/emicklei/go-restful/log"
)

// Instance returns a default logger.
func Instance() *logrus.Logger {
	return logrus.StandardLogger()
}

func SetupAndHarmonize(verbose bool) {
	restfullog.SetLogger(new(DebugLogger))
	if verbose {
		Instance().Info("Turn on verbose logging")
		logrus.SetLevel(logrus.DebugLevel)
		restful.TraceLogger(&DebugLogger{Prefix: "[restful/trace]"})
		restful.EnableTracing(true)
	}
}
