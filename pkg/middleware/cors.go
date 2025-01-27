package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/even44/JobsearchAPI/pkg/initializers"
)

var corsLogger *log.Logger = log.New(os.Stdout, "[CORS] ", log.Ldate+log.Ltime+log.Lmsgprefix)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Origin") == initializers.ApiTrustedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		} else {
			corsLogger.Printf("Received request with wrong origin, %s, from: %s", r.Header.Get("Origin"), r.Host)
			http.Error(w, "bad origin", http.StatusUnauthorized)
			return
		}

		corsLogger.Println("Adding CORS Headers")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")
		next.ServeHTTP(w, r)
	})
}
