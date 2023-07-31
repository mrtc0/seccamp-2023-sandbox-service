package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func payController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "Payment completion! Thank you ~~~ 💸")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})
}

func healthCheckController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
}

func main() {
	logger := log.New(os.Stdout, "payments: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", healthCheckController())
	router.Handle("/pay", payController())

	server := &http.Server{
		Addr:     ":7000",
		Handler:  logging(logger)(router),
		ErrorLog: logger,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", ":7000", err)
	}
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID := r.Header.Get("X-Request-Id")
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
