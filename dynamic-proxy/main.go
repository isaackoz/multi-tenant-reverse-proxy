package main

import (
	"context"
	mycache "dynamic-proxy/cache"
	"dynamic-proxy/config"
	"dynamic-proxy/db"
	httphelper "dynamic-proxy/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/joho/godotenv"
	// "time"
)

func NewServer(logger *log.Logger, config *config.Config, rateLimitCache mycache.RateLimitCache, tenantCache mycache.TenantCache) http.Handler {

	// Create a router
	r := httphelper.NewRouter(loggingMiddleware(logger))

	// Handle proxy
	r.Any("/ask", dynamicProxyHandler(logger, config, rateLimitCache, tenantCache), rateLimitByIpMiddleware(logger, rateLimitCache))

	// Health check
	r.Get("/healthz", handleHealthz(logger))

	// Handle clearing cache for specified tenant
	r.Delete("/invalidate", handleInvalidate(logger, config, tenantCache))

	return r
}

func dynamicProxyHandler(logger *log.Logger, config *config.Config, rateLimitCache mycache.RateLimitCache, tenantCache mycache.TenantCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the domain from the query param /?domain=...
		domain := r.URL.Query().Get("domain")

		if domain == "" {
			http.Error(w, "Invalid or missing domain", http.StatusBadRequest)
			return
		}

		isAllowed, err := determineBackend(domain, tenantCache)

		if err != nil {
			handleNotFoundRateLimiter(r.RemoteAddr, logger, rateLimitCache)
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			logger.Printf("Error determining backend for host %s: %v\n", domain, err)
			return
		}

		if isAllowed == false {
			http.Error(w, "Not allowed", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func handleHealthz(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Health check")
		w.WriteHeader(http.StatusOK)
	}
}

func handleInvalidate(logger *log.Logger, config *config.Config, tenantCache mycache.TenantCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get auth bearer token
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		bearerToken := authHeader[1]

		// Check if token is valid
		if bearerToken != config.AuthToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get the hostname from the query
		hostname := r.URL.Query().Get("hostname")

		if hostname == "" {
			http.Error(w, "Invalid or missing hostname", http.StatusBadRequest)
			return
		}

		// Finally, delete the tenant id in the cache
		err := tenantCache.Delete(context.TODO(), fmt.Sprintf("t-%s", hostname))

		if err != nil {
			http.Error(w, "Error deleting tenant", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func rateLimitByIpMiddleware(logger *log.Logger, rateLimitCache mycache.RateLimitCache) httphelper.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("Request from IP: %s\n", r.RemoteAddr)
			ctx := context.TODO()
			// Check if user is rate limited
			val, err := rateLimitCache.Get(ctx, fmt.Sprintf("rl-%s", r.RemoteAddr))
			if err != nil {
				// If we do not get a value in the store, we can assume the user is not being rate limited
				next.ServeHTTP(w, r)
				return
			}
			logger.Printf("Rate limit value: %v\n", val)
			if val >= 10 {
				http.Error(w, "Rate limited", http.StatusTooManyRequests)
				logger.Printf("Rate limited user: %s\n", r.RemoteAddr)
				return
			} else {
				next.ServeHTTP(w, r)
				return
			}
		})
	}
}

func handleNotFoundRateLimiter(remoteAddr string, logger *log.Logger, rateLimitCache mycache.RateLimitCache) error {
	ctx := context.TODO()
	// If a requested route is not found, we will add 1 to the rate limiter
	currentLimit := 0
	// Get their current val
	currVal, err := rateLimitCache.Get(ctx, fmt.Sprintf("rl-%s", remoteAddr))
	if err != nil {
		// Key doesn't exist. That's ok, we'll just increase by 1
		currentLimit = 1
	} else {
		// add int32
		currentLimit = currVal + 1
	}

	// Update the rate limit
	return rateLimitCache.Set(ctx, fmt.Sprintf("rl-%s", remoteAddr), currentLimit, store.WithCost(1), store.WithExpiration(1*time.Minute))
}

func loggingMiddleware(logger *log.Logger) httphelper.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		http.Error(w, "Proxy Error", http.StatusBadGateway)
	}
}

// determineBackend
func determineBackend(hostname string, tenantCache mycache.TenantCache) (bool, error) {
	// TODO: implement cache/db lookup
	ctx := context.TODO()

	fmt.Println("Determining backend for host: ", hostname)

	val, err := tenantCache.Get(ctx, fmt.Sprintf("t-%s", hostname))

	if err != nil {
		return false, errors.New("Tenant does not exist")
	}

	if val == "" {
		return false, errors.New("Error determining backend")
	}

	return true, nil
}

func run(
	ctx context.Context,
	logger *log.Logger,
	config *config.Config,
	rateLimitCache mycache.RateLimitCache,
	tenantCache mycache.TenantCache,
) error {
	srv := NewServer(logger, config, rateLimitCache, tenantCache)
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
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s\n", err)
	}

	envToInt := func(path string) int {
		val, err := strconv.Atoi(os.Getenv(path))
		if err != nil {
			logger.Fatalf("Error converting %s to int: %s\n", path, err)
		}
		logger.Println(path, ":", val)
		return val
	}

	envToStr := func(path string) string {
		val := os.Getenv(path)
		if val == "" {
			logger.Fatalf("Error converting %s to string: empty value\n", path)
		}
		logger.Println(path, ":", val)
		return val
	}

	config := &config.Config{
		Host:      envToStr("HOST"),
		Port:      envToStr("PORT"),
		RedisAddr: envToStr("REDIS_URI"),
		PostgresAddr: config.PostgresAddr{
			Host:     envToStr("PG_HOST"),
			Port:     envToInt("PG_PORT"),
			User:     envToStr("PG_USER"),
			Password: envToStr("PG_PASS"),
			Dbname:   envToStr("PG_DB_NAME"),
		},
		AuthToken: envToStr("AUTH_TOKEN"),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup cache
	rateLimitCache := mycache.CreateRateLimitCacheManager()
	tenantCache := mycache.CreateTenantCacheManager(config)

	// Setup postgres
	err = db.InitPostgres(config.PostgresAddr)

	if err != nil {
		logger.Fatalf("Error connecting to Postgres: %s\n", err)
		panic(err)
	}
	defer db.PostgresDB.Close()

	if err := run(ctx, logger, config, rateLimitCache, tenantCache); err != nil {
		logger.Fatalf("Server exited with error: %s\n", err)
	}
}
