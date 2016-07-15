// Package log provides logging for the kube-volume-freezer project.
package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	restfullog "github.com/emicklei/go-restful/log"
)

// Instance returns the default logger.
func Instance() *logrus.Logger {
	return logrus.StandardLogger()
}

// SetupAndHarmonize configures logging including third-party loggers.
func SetupAndHarmonize(verbose bool) {
	restfullog.SetLogger(new(DebugLogger))
	if verbose {
		Instance().Info("Turn on verbose logging")
		logrus.SetLevel(logrus.DebugLevel)
		restful.TraceLogger(&DebugLogger{Prefix: "[restful/trace]"})
		restful.EnableTracing(true)
	}
}
