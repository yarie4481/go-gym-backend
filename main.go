package main

import (
	"context"
	"fmt"
	"go-blog/controllers"
	"go-blog/internal/config"
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"go-blog/logger"
	"go-blog/middlewares"
	"go-blog/repositories"
	"go-blog/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// Application represents the main application structure
type Application struct {
	Config    *gorm.DB
	Logger    *logrus.Logger
	Router    *gin.Engine
	Server    *http.Server
	Health    *HealthChecker
	DB        *gorm.DB
}

// HealthChecker manages health checks
type HealthChecker struct {
	checks map[string]func() error
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]func() error),
	}
}

func (h *HealthChecker) AddCheck(name string, check func() error) {
	h.checks[name] = check
}

func (h *HealthChecker) IsReady() bool {
	for name, check := range h.checks {
		if err := check(); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"check": name,
				"error": err,
			}).Error("Health check failed")
			return false
		}
	}
	return true
}

// Build information (set during build)
var (
	version    = "dev"
	commitHash = "unknown"
	buildDate  = "unknown"
)

// --- Custom Errors/API Response Structure ---
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error
}

func (e *APIError) Error() string {
	return e.Message
}

// --- Metrics Initialization ---
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	appVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_version",
			Help: "Application version information",
		},
		[]string{"version", "commit", "build_date"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration, appVersion)
}

// --- Middlewares ---

// PrometheusMiddleware collects metrics and performs enhanced logging
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		httpRequests.WithLabelValues(method, path, http.StatusText(status)).Inc()
		httpDuration.WithLabelValues(method, path, http.StatusText(status)).Observe(duration)

		logger.Log.WithFields(logrus.Fields{
			"requestId": c.GetString("requestID"),
			"method":    method,
			"path":      path,
			"status":    status,
			"duration":  duration,
			"clientIP":  c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}).Info("HTTP request processed")
	}
}

// RequestIDMiddleware injects a unique ID into the context for tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// ErrorHandlingMiddleware catches errors and standardizes responses
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err, exists := c.Get("error"); exists {
			var apiErr *APIError
			if customErr, ok := err.(*APIError); ok {
				apiErr = customErr
			} else {
				apiErr = &APIError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
					Details: "An unexpected error occurred.",
					Err:     err.(error),
				}
			}

			logger.Log.WithFields(logrus.Fields{
				"requestId": c.GetString("requestID"),
				"status":    apiErr.Code,
				"error":     apiErr.Err.Error(),
			}).Error("API request failed with error")

			c.AbortWithStatusJSON(apiErr.Code, gin.H{
				"code":    apiErr.Code,
				"message": apiErr.Message,
				"details": apiErr.Details,
			})
		}
	}
}

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// --- Handlers ---

// HealthHandler returns service health status
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   version,
	})
}

// ReadyHandler returns service readiness status
func ReadyHandler(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !app.Health.IsReady() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not ready",
				"message": "Service is not ready to handle requests",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"message": "Service is ready to handle requests",
		})
	}
}

// VersionHandler returns build information
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    version,
		"commit":     commitHash,
		"build_date": buildDate,
	})
}

// CBEBirrTestHandler - Simple endpoint to test token capture
func CBEBirrTestHandler(c *gin.Context) {
	// Log everything about the request
	logger.Log.Info("\n=== CBEBirr TEST ENDPOINT HIT ===")
	logger.Log.WithFields(logrus.Fields{
		"time":      time.Now().Format(time.RFC3339),
		"url":       c.Request.URL.String(),
		"method":    c.Request.Method,
		"remote_ip": c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
	}).Info("Request details")

	// Log all headers
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	logger.Log.WithField("headers", headers).Info("Request headers")

	// Check Authorization header
	authHeader := c.GetHeader("Authorization")
	logger.Log.WithField("authorization_header", authHeader).Info("Authorization header")

	if authHeader == "" {
		logger.Log.Warn("❌ NO Authorization header found!")
		c.JSON(200, gin.H{
			"success": false,
			"message": "NO Authorization header found",
			"note":    "CBEBirr should send: Authorization: Bearer YOUR_TOKEN",
			"test_by": "Yared from Adiamat",
		})
		return
	}

	// Check if it starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		logger.Log.WithFields(logrus.Fields{
			"expected": "Bearer {token}",
			"got":      authHeader,
		}).Warn("❌ Authorization header format wrong!")
		
		c.JSON(200, gin.H{
			"success":         false,
			"message":         "Authorization header should start with 'Bearer '",
			"received_header": authHeader,
			"test_by":         "Yared from Adiamat",
		})
		return
	}

	// Extract token
	token := authHeader[7:]
	logger.Log.WithFields(logrus.Fields{
		"token":        token,
		"token_length": len(token),
	}).Info("✅ TOKEN EXTRACTED SUCCESSFULLY!")

	// Return success response
	c.JSON(200, gin.H{
		"success":            true,
		"message":            "Token captured successfully!",
		"token_length":       len(token),
		"token_first_10":     getFirstChars(token, 10),
		"token_last_10":      getLastChars(token, 10),
		"note":               "Token received correctly. You can proceed with full integration.",
		"test_by":            "Yared from Adiamat",
		"test_endpoint":      "http://54.159.60.214:8802/api/cbebirr-test",
		"timestamp":          time.Now().Format(time.RFC3339),
	})
}

// Helper functions for token display
func getFirstChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func getLastChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// CBEBirrCallbackHandler - For POST callback testing
func CBEBirrCallbackHandler(c *gin.Context) {
	logger.Log.Info("\n=== CBEBirr CALLBACK TEST ===")
	
	// Log headers
	authHeader := c.GetHeader("Authorization")
	logger.Log.WithField("authorization", authHeader).Info("Callback headers")
	
	// Try to parse body
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err == nil {
		logger.Log.WithField("body", body).Info("Callback body")
	}
	
	// Simple response
	c.JSON(200, gin.H{
		"success":           true,
		"message":           "Callback received",
		"has_auth_header":   authHeader != "",
		"auth_header_exists": len(authHeader) > 0,
		"test_by":           "Yared from Adiamat",
	})
}

// --- Application Initialization ---

func initializeApplication(ctx context.Context) (*Application, error) {
	// 1. Initialize logger
	logger.Init()
	logger.Log.Info("Starting Go-Blog server...")

	// 2. Load configuration
	config.InitDB()
	config.InitJWT()
	
	// 3. Initialize database
	// Configure connection pool
	sqlDB, err := config.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}
	
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// 4. Run migrations
	models.MigrateModels(config.DB)

	// 5. Setup health checks
	healthChecker := NewHealthChecker()
	healthChecker.AddCheck("database", func() error {
		return sqlDB.Ping()
	})

	// 6. Initialize Gin
	router := setupRouter(healthChecker)

	// 7. Create HTTP server
	server := &http.Server{
		Addr:         ":8888",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. Set version metrics
	appVersion.WithLabelValues(version, commitHash, buildDate).Set(1)

	return &Application{
		Config: config.DB,
		Logger: logger.Log,
		Router: router,
		Server: server,
		Health: healthChecker,
		DB:     config.DB,
	}, nil
}

func setupRouter(healthChecker *HealthChecker) *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = false

	// CORS Middleware - SIMPLE WILDCARD VERSION
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Use wildcard
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders: []string{"Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: false, // IMPORTANT: Must be false for wildcard *
		MaxAge: 12 * time.Hour,
	}))

	// Debug middleware to log all requests
	r.Use(func(c *gin.Context) {
		logger.Log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"full_path": c.FullPath(),
		}).Info("Incoming request")
		c.Next()
	})

	// Global Middlewares
	r.Use(RequestIDMiddleware())
	r.Use(ErrorHandlingMiddleware())
	r.Use(TimeoutMiddleware(30 * time.Second))
	r.Use(PrometheusMiddleware())

	// Basic routes
	r.GET("/health", HealthHandler)
	r.GET("/ready", ReadyHandler(&Application{Health: healthChecker}))
	r.GET("/version", VersionHandler)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ✅ ADD CBEBirr Test Endpoints
	r.GET("/api/cbebirr-test", CBEBirrTestHandler)
	r.POST("/api/cbebirr-test", CBEBirrCallbackHandler)
	r.GET("/cbebirr-test", CBEBirrTestHandler) // Alternative endpoint

	// OPTIONS handler for preflight requests - SIMPLIFIED
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Status(http.StatusOK)
	})

	// Initialize dependencies and register routes
	registerAllRoutes(r)

	// Print all registered routes for debugging
	printRoutes(r)

	return r
}

// Add this function to debug routes
func printRoutes(r *gin.Engine) {
	routes := r.Routes()
	logger.Log.Info("=== Registered Routes ===")
	for _, route := range routes {
		logger.Log.WithFields(logrus.Fields{
			"method": route.Method,
			"path":   route.Path,
		}).Info("Route")
	}
	logger.Log.Info("=== End Registered Routes ===")
}

func registerAllRoutes(r *gin.Engine) {
	// Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepo, config.DB)
	authController := controllers.NewAuthController(authService)

	attendanceRepo := repositories.NewAttendanceRepository(config.DB)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	attendanceController := controllers.NewAttendanceController(attendanceService)

	classSessionRepo := repositories.NewClassSessionRepository(config.DB)
	classSessionService := services.NewClassSessionService(classSessionRepo)
	classSessionController := controllers.NewClassSessionController(classSessionService)

	classRepo := repositories.NewClassRepository(config.DB)
	classService := services.NewClassService(classRepo)
	classController := controllers.NewClassController(classService)

	gymRepo := repositories.NewGymRepository(config.DB)
	gymService := services.NewGymService(gymRepo)
	gymController := controllers.NewGymController(gymService)

	planRepo := repositories.NewPlanRepository(config.DB)
	memberRepo := repositories.NewMembershipRepository(config.DB)
	paymentRepo := repositories.NewPaymentRepository(config.DB)

	// Services
	planService := services.NewPlanService(planRepo)
	memberService := services.NewMembershipService(memberRepo)
	paymentService := services.NewPaymentService(paymentRepo)

	// Controllers
	planController := controllers.NewPlanController(planService)
	memberController := controllers.NewMembershipController(memberService)
	paymentController := controllers.NewPaymentController(paymentService)

	memberRepo1 := repositories.NewMemberRepository(config.DB)
	member1Service := services.NewMemberService(memberRepo1)
	memberController1 := controllers.NewMemberController(member1Service)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/alluser", authController.GetAllUsers)
		auth.GET("/trainer", authController.GetTrainers)
	}

	// Register all routes
	routes.RegisterAttendanceRoutes(r, attendanceController)
	routes.RegisterClassSessionRoutes(r, classSessionController)
	routes.RegisterClassRoutes(r, classController)
	routes.RegisterGymRoutes(r, gymController)
	routes.RegisterRoutes(r, planController, memberController, paymentController)
	routes.RegisterMemberRoutes(r, memberController1)

	// Protected routes
	protected := r.Group("/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			requestID := c.GetString("requestID")
			c.JSON(http.StatusOK, gin.H{
				"message":   "This is protected",
				"request_id": requestID,
			})
		})
	}
}

// --- Graceful Shutdown ---

func setupGracefulShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		logger.Log.Infof("Received signal: %v", sig)
		cancel()
	}()
}

func (app *Application) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		app.Logger.Infof("Starting server on %s", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server failed: %w", err)
		}
		return nil
	})

	// Graceful shutdown handler
	g.Go(func() error {
		<-ctx.Done()
		
		app.Logger.Info("Shutting down server gracefully...")
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := app.Server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		// Close database connection
		if sqlDB, err := app.DB.DB(); err == nil {
			sqlDB.Close()
		}

		app.Logger.Info("Server shutdown complete")
		return nil
	})

	return g.Wait()
}

// --- Main Function ---

func main() {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	setupGracefulShutdown(cancel)

	// Initialize application
	app, err := initializeApplication(ctx)
	if err != nil {
		logger.Log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the application
	if err := app.Run(ctx); err != nil {
		app.Logger.Fatalf("Application failed: %v", err)
	}
}
