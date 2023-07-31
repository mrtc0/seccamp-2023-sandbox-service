package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const timeForHeavyProcessing = 20 * time.Second

type Item struct {
	ID   int
	Name string
}

var itemsMockResponse = []Item{
	{ID: 1, Name: "Item 1"},
	{ID: 2, Name: "Item 2"},
	{ID: 3, Name: "Item 3"},
	{ID: 4, Name: "Item 4"},
	{ID: 5, Name: "Item 5"},
}

func itemsController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSlow := r.URL.Query().Get("slow")

		if isSlow == "true" {
			time.Sleep(timeForHeavyProcessing)
		}

		data, _ := json.Marshal(itemsMockResponse)
		fmt.Fprintf(w, "%s", data)
		w.Header().Set("Content-Type", "application/json")
	})
}

func healthCheckController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
}

func main() {
	logger := log.New(os.Stdout, "item-api: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", healthCheckController())
	router.Handle("/items", itemsController())

	server := &http.Server{
		Addr:     ":9000",
		Handler:  logging(logger)(router),
		ErrorLog: logger,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", ":9000", err)
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
