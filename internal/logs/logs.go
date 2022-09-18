// Package logs facilitates creating of configured Logrus logger.
// Logger created with NewLogger() should be preferred over global Logger.
// Use logrus.WithField() and logrus.WithFields() methods with "dash-case" keys for additional log parameters.
package logs

import "github.com/sirupsen/logrus"

// TODO(dfurman): make configurable: level, formatter.
// TODO(dfurman): return Logger interface to decouple logger users from Logrus

const (
	defaultLevel = logrus.DebugLevel
)

// Configure configures global Logrus logger.
func Configure() {
	logrus.SetLevel(defaultLevel)
}

func NewLogger(name string) *logrus.Entry {
	l := logrus.New()
	l.Level = defaultLevel
	return l.WithField("logger", name)
}
