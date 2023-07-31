package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type key int

const (
	requestIDKey key = 0
)

const defaultTimeout = 5 * time.Second

const itemListPath = "/items"
const payPath = "/pay"

func httpClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimeout,
	}
}

func buildURL(endpoint, path string) (*url.URL, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = filepath.Join(u.Path, path)
	return u, nil
}

type ItemService struct {
	URL    string
	Client *http.Client
}

type PaymentService struct {
	URL    string
	Client *http.Client
}

func NewItemService(isSlow string) (*ItemService, error) {
	u, err := buildURL(os.Getenv("ITEMS_API_ADDR"), itemListPath)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("slow", isSlow)
	u.RawQuery = q.Encode()

	return &ItemService{
		URL:    u.String(),
		Client: httpClient(),
	}, nil
}

func (s *ItemService) FetchItems(ctx context.Context) ([]byte, error) {
	req, _ := http.NewRequest("GET", s.URL, nil)
	req.Header.Set("X-Request-Id", ctx.Value(requestIDKey).(string))

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewPaymentService() (*PaymentService, error) {
	u, err := buildURL(os.Getenv("PAYMENTS_API_ADDR"), payPath)
	if err != nil {
		return nil, err
	}

	return &PaymentService{
		URL:    u.String(),
		Client: httpClient(),
	}, nil
}

func (s *PaymentService) Pay(ctx context.Context) ([]byte, error) {
	req, _ := http.NewRequest("GET", s.URL, nil)
	req.Header.Set("X-Request-Id", ctx.Value(requestIDKey).(string))

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func fetchItemsController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSlow := r.URL.Query().Get("slow")

		svc, err := NewItemService(isSlow)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		items, err := svc.FetchItems(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(items)
	})
}

func uploadImageController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// mock
	})
}

func paymentController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svc, err := NewPaymentService()
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		result, err := svc.Pay(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})
}

func healthCheckController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func main() {
	logger := log.New(os.Stdout, "backend: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", healthCheckController())
	router.Handle("/items", fetchItemsController())
	router.Handle("/upload", uploadImageController())
	router.Handle("/payment", paymentController())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:     ":8000",
		Handler:  tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog: logger,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", ":8000", err)
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
