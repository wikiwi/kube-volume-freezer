/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package log

import (
	"github.com/Sirupsen/logrus"
)

// DebugLogger implements the go-restful logger printing to debug using logrus.
type DebugLogger struct {
	// Optionally prepend a prefix on each print.
	Prefix string
}

// Print calls logrus.Debug().
func (l *DebugLogger) Print(v ...interface{}) {
	if l.Prefix != "" {
		v = append([]interface{}{l.Prefix + " "}, v...)
	}
	logrus.Debug(v...)
}

// Printf calls logrus.Debuf().
func (l *DebugLogger) Printf(format string, v ...interface{}) {
	if l.Prefix != "" {
		format = l.Prefix + " " + format
	}
	// Logrus adds a \n at the end of the string unlike the standard logger.
	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
	}
	logrus.Debugf(format, v...)
}
