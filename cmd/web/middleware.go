package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/justinas/nosurf"
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
func (app *application) logRequest(next http.Handler) http.Handler {

	// Open the log file
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, _ = fmt.Fprintln(f, "-----New Log Session-----")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("IP: %s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

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

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication redirects an unauthenticated user to login before performing specific actions
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set the "Cache-Control: no-store" so pages that require authentication are not stored in users browser cache
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// noSurf protects handlers or forms from CSRF (Cross-site Request Forgery)
func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})

	return csrfHandler
}