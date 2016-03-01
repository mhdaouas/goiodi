package main

import (
	"github.com/op/go-logging"
	"net/http"
	"os"
	"time"
)

// Application prefix used in syslog
const logPrefix = "goiodi"

var Log = logging.MustGetLogger(logPrefix)

var stderrFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} [%{level:.4s}] %{id:03x} ▶ %{shortfile} %{shortfunc}%{color:reset} %{message}",
)
var syslogFormat = logging.MustStringFormatter(
	"%{level:.1s}%{id:03x} %{shortfile} %{shortfunc} ▶ %{message}",
)

func InitLogger() {
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	syslogBackend, err := logging.NewSyslogBackend(logPrefix)
	if err != nil {
		Log.Fatal(err)
	}

	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stderrFormat)
	syslogFormatter := logging.NewBackendFormatter(syslogBackend, syslogFormat)

	// Only errors and more severe messages should be sent to the syslog
	syslogLeveled := logging.AddModuleLevel(syslogFormatter)
	syslogLeveled.SetLevel(logging.WARNING, "")

	// Set the backends to be used.
	logging.SetBackend(stderrFormatter, syslogLeveled)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		Log.Info(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
