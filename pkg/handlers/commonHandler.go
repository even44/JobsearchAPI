package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/even44/JobsearchAPI/pkg/initializers"
)

var logger *log.Logger = log.New(os.Stdout, "[COMMON HANDLER]", log.Ldate+log.Ltime+log.Lmsgprefix)

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func InvalidEmailOrPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Wrong email or password"))
}

func PreFlightHandler(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}
func checkOrigin(w *http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Origin") == initializers.ApiTrustedOrigin {
		(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		return true
	} else {
		logger.Printf("Received request with wrong origin, %s, from: %s", r.Header.Get("Origin"), r.Host)
		InternalServerErrorHandler((*w), r)
		return false
	}
}
func enableCors(w *http.ResponseWriter) {

	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")

}
