package main

import (
	"context"
	"fmt"
	"log"

	"net/http"

	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/a-dev-mobile/kidneysmart-auth/docs"
	"github.com/a-dev-mobile/kidneysmart-auth/database/mongo"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/pkg/emailclient"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/register"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/verifycode"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/logging"

	"golang.org/x/exp/slog"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/middleware"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

)

func main() {
	cfg, lg := initializeApp()

	// Set up the database connection
	db, cleanup := setupDatabase(cfg, lg)
	defer cleanup()
	setGinMode(cfg)

	// Initialize gRPC connection to SMTP server with updated security settings
	smtpConn, err := grpc.Dial(
		cfg.ExternalService.SmtpServer.Grpc.Host+":"+cfg.ExternalService.SmtpServer.Grpc.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),    // Updated to use WithTransportCredentials
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(5<<20)), // Set max message size if needed
	)
	if err != nil {
		lg.Error("Failed to connect to SMTP server:", logging.Err(err))
		os.Exit(1)
	}
	defer smtpConn.Close()

	emailClient := emailclient.NewEmailClient(smtpConn,lg)

	// Set up your server's routes and handlers
	router := setupRouter(cfg, lg)

	// Create a new context for the Register handler including the email client
	hctxRegister := register.NewRegisterServiceContext(db, lg, cfg, emailClient)
	router.POST("kidneysmart-auth/v1/register", hctxRegister.RegisterUserHandler)
	//
	hctxVerifyCode := verifycode.NewVerifyCodeServiceContext(db, lg, cfg)
	router.POST("kidneysmart-auth/v1/verify-code", hctxVerifyCode.VerifyCodeHandler)
	
	lg.Info("Environment used", ".env", cfg.Environment)
	// lg.Debug("Rest Server starting", "config_json", cfg)
	// docs
	docs.SwaggerInfo.Title = "Kidneysmart Auth API"
	docs.SwaggerInfo.Description = "Authentication in the KidneySmart application"
	docs.SwaggerInfo.Host = "wayofdt.com"
	
	docs.SwaggerInfo.BasePath = "/kidneysmart-auth/v1/"
	docs.SwaggerInfo.Version = "v1"
	router.GET("/kidneysmart-auth/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	startServer(cfg, router, lg)
}

// initializeApp sets up the application environment, configuration, and logger.
// It determines the application's running environment, loads the appropriate configuration,
// and initializes the logging system.
func initializeApp() (*config.Config, *slog.Logger) {

	cfg := getConfigOrFail()

	lg := logging.SetupLogger(cfg)

	return cfg, lg
}

func getConfigOrFail() *config.Config {

	cfg, err := config.LoadConfig("../config", "config.yaml")

	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	return cfg
}

// setGinMode configures the Gin mode (debug or release) based on the application's configuration.
func setGinMode(cfg *config.Config) {
	switch cfg.ClientConnection.GinMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// setupDatabase initializes the MongoDB database connection using the provided configuration and logger.
// It returns a MongoDB client and a cleanup function to disconnect from the database.
// setupDatabase initializes the MongoDB database connection using the provided configuration and logger.
// It returns a MongoDB client and a cleanup function to disconnect from the database.
func setupDatabase(cfg *config.Config, lg *slog.Logger) (*mongodriver.Client, func()) {
	dbUser := cfg.Database.User
	dbPassword := cfg.Database.Password
	dbPort := cfg.Database.Port
	dbHost := cfg.Database.Host

	db, err := mongo.GetDB(context.Background(), dbUser, dbPassword, dbHost, dbPort)
	if err != nil {
		lg.Error("Error initializing database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return db, func() {
		if err := db.Disconnect(context.Background()); err != nil {
			lg.Error("Error disconnecting from database", "error", err.Error())
		}
	}
}

// setupRouter initializes and returns a new Gin router configured with middleware and routes.
func setupRouter(cfg *config.Config, lg *slog.Logger) *gin.Engine {
	// Create a new router
	router := gin.New()
	// Apply global middleware
	router.Use(gin.Recovery()) // Recovery middleware от Gin
	router.Use(gin.Logger())   // Logging middleware от Gin
	router.Use(middleware.CORSMiddleware(*cfg, lg))
	router.Use(middleware.TrustProxyHeader())

	// Adding custom middleware to recover from a panic
	router.Use(middleware.RecoveryMiddleware(lg))

	return router
}

// startServer starts the HTTP server using the provided configuration and router,
// and handles graceful shutdown on receiving quit signals.
func startServer(cfg *config.Config, router *gin.Engine, lg *slog.Logger) {
	serverAddr := fmt.Sprintf(":%s", cfg.ClientConnection.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Error("Error running server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	gracefulShutdown(srv, lg)
}
func gracefulShutdown(srv *http.Server, lg *slog.Logger) {
	quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt) // Uncomment this line if you want to listen to OS interrupt signals

	<-quit
	lg.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lg.Error("Server forced to shutdown:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	lg.Info("Server exiting")
}
func LogHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentTime := time.Now().Format("2006/01/02 15:04:05")
		log.Printf("\n[%s] Headers for %s %s:\n", currentTime, c.Request.Method, c.Request.URL.Path)
		for k, v := range c.Request.Header {
			log.Printf("  %s: %s", k, strings.Join(v, ","))
		}
		c.Next()
	}
}
