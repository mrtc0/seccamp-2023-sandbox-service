package main

import (
	"context"
	"fmt"
	"io"
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

// proxy's port
const proxyPort = 8080

// payment's container exposed port
const paymentServicePort = 7000

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res, err := p.forwardRequest(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	p.writeResponse(w, res)
}

func isAllowedRequest(req *http.Request) bool {
	// TODO: X-Internal-Token が設定されていれば true を返し、設定されていなければ false を返す
	return true
}

// 受け取ったリクエストを paymentServicePort に転送する
func (p *Proxy) forwardRequest(req *http.Request) (*http.Response, error) {
	if !isAllowedRequest(req) {
		return nil, fmt.Errorf("not allowed request")
	}

	proxyUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", paymentServicePort))
	if err != nil {
		return nil, err
	}

	proxyUrl.Path = filepath.Join(proxyUrl.Path, req.RequestURI)

	httpClient := http.Client{}
	proxyReq, err := http.NewRequest(req.Method, proxyUrl.String(), req.Body)
	proxyReq.Header.Set("X-Request-Id", req.Context().Value(requestIDKey).(string))

	res, err := httpClient.Do(proxyReq)

	return res, err
}

// payment から受け取ったレスポンスを返す
func (p *Proxy) writeResponse(w http.ResponseWriter, res *http.Response) {
	// Copy response headers
	for name, values := range res.Header {
		w.Header()[name] = values
	}

	w.Header().Set("Server", "seccamp-proxy")
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func main() {
	logger := log.New(os.Stdout, "sidecar: ", log.LstdFlags)
	logger.Println("sidecar proxy is starting...")

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	proxy := &Proxy{}

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", proxyPort),
		Handler:  tracing(nextRequestID)(logging(logger)(proxy)),
		ErrorLog: logger,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", fmt.Sprintf(":%d", proxyPort), err)
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
