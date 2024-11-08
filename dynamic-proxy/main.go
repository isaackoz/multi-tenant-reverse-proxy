package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	// "time"
)

func NewServer(logger *log.Logger, config *Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger, config)

	var handler http.Handler = mux
	handler = loggingMiddleware(logger, handler)
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	config *Config,
) {
	mux.HandleFunc("/healthz", handleHealthz(logger))
	mux.HandleFunc("/", dynamicProxyHandler(logger, config))
}

func dynamicProxyHandler(logger *log.Logger, config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHttps := false
		if r.TLS != nil {
			isHttps = true
		}

		// startTime := time.Now()
		targetHostname, err := determineBackend(r.Host, config.TargetBackend, isHttps)

		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			logger.Printf("Error determining backend for host %s: %v\n", r.Host, err)
		}

		proxy, err := NewProxy(config.TargetBackend, targetHostname)
		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			logger.Printf("Error creating proxy for target %s: %v\n", targetHostname, err)
		}
		proxy.ServeHTTP(w, r)
		// duration := time.Now().Sub(startTime)
		// logger.Printf("Request took %s\n", duration)
	}
}

func handleHealthz(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Health check")
		w.WriteHeader(http.StatusOK)
	}
}

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NewProxy(targetBackend string, targetHostname string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHostname)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	// proxy := &httputil.ReverseProxy{
	// 	Rewrite: func(r *httputil.ProxyRequest) {
	// 		r.SetURL(backendUrl)
	// 	},
	// }
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req, url.Hostname())
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request, targetHostname string) {
	// TODO: interact with cache/db to modify upstream
	req.Host = targetHostname
	req.Header.Set("Host", targetHostname)
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		// Example modification; in a real setup, customize as needed
		if resp.StatusCode == http.StatusBadGateway {
			return errors.New("response body is invalid")
		}
		return nil
	}
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		http.Error(w, "Proxy Error", http.StatusBadGateway)
	}
}

// determineBackend
func determineBackend(hostname string, target string, isHttps bool) (string, error) {
	// TODO: implement cache/db lookup

	// Check KV cache for host
	// i.e. isaac.com will return an arbitrary tenant id such as isaac
	tenantId := ""
	if hostname == "1saac.com:8080" {
		tenantId = "isaac"
	}

	// Todo handle error
	if tenantId == "" {
		return "", errors.New("Error determining backend")
	}

	// If isHttps then add https
	if isHttps {
		return fmt.Sprintf("%s%s.%s", "https://", tenantId, target), nil
	} else {
		return fmt.Sprintf("%s%s.%s", "http://", tenantId, target), nil
	}
}

type Config struct {
	Host          string
	Port          string
	TargetBackend string
}

func run(
	ctx context.Context,
	logger *log.Logger,
	config *Config,
) error {
	srv := NewServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	go func() {
		logger.Printf("Listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error listening and server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Printf("Error shutting down HTTP server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	config := &Config{
		Host:          "0.0.0.0",
		Port:          "8080",
		TargetBackend: "localhost.test:5173",
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, logger, config); err != nil {
		logger.Fatalf("Server exited with error: %s\n", err)
	}
}
