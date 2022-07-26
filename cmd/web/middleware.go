package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var current_time = time.Now().Local().Format("01/02/2006 15:04:05")

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nostiff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// --log: Check the log of each request
func (log *application) logRequest(next http.Handler) http.Handler {

	// Open the log file
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, _ = fmt.Fprintln(f, "-----New Log Session-----")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.infoLog.Printf("IP: %s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		// newLog := log.New(os.Stdout, "INFO:\t", dateLog.Ldate|dateLog.Ltime)

		// Append to the log file

		logToFile := fmt.Sprintf("Time: %v, IP: %s - %s %s %s", current_time, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		_, err = fmt.Fprintln(f, logToFile)
		if err != nil {
			fmt.Println(err)
		}

		next.ServeHTTP(w, r)
	})
}

func (log *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				log.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
