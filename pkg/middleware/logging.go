package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var logLogger *log.Logger = log.New(os.Stdout, "[REQUEST] ", log.Ldate+log.Ltime+log.Lmsgprefix)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logLogger.Println(r.Method, r.URL.Path, time.Since(start))
	})
}
