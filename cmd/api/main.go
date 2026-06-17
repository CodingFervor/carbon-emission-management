package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/cache"
	"github.com/CodingFervor/carbon-emission-management/internal/config"
	"github.com/CodingFervor/carbon-emission-management/internal/database"
	"github.com/CodingFervor/carbon-emission-management/internal/handler"
	"github.com/CodingFervor/carbon-emission-management/internal/repository"
	"github.com/CodingFervor/carbon-emission-management/internal/server"
	"github.com/CodingFervor/carbon-emission-management/pkg/jwt"
	"github.com/CodingFervor/carbon-emission-management/pkg/logger"
)

func main() {
	// Load configuration (falls back to sensible defaults if file is absent).
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Warn("config file not loaded, using defaults", "error", err)
		cfg = defaultConfig()
	}
	gin.SetMode(cfg.Server.Mode)
	logger.SetLevel(cfg.Server.Mode)
	jwt.SetSecret(cfg.JWT.Secret)
	// Connect infrastructure. Failures are logged but do not abort startup so
	// the API can still serve health/liveness probes.
	if err := database.Connect(cfg.Database); err != nil {
		logger.Error("failed to connect database", "error", err)
	} else {
		defer database.Close()
	}
	if err := cache.Connect(cfg.Redis); err != nil {
		logger.Error("failed to connect redis", "error", err)
	} else {
		defer cache.Close()
	}

	// Build the handler graph, injecting all repositories.
	h := newHandler()

	r := server.New(h)

	addr := ":" + strconv.Itoa(cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 10 * time.Second}

	go func() {
		logger.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on interrupt / terminate signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}
	logger.Info("server exited")
}

// newHandler wires every repository into a single handler container.
// Repositories degrade gracefully when the DB is unavailable: a nil *sql.DB
// is passed in that case so the app still boots (endpoints will return 500).
func newHandler() *handler.Handler {
	db := database.DB
	return &handler.Handler{
		Org:        repository.NewOrganizationRepo(db),
		Facility:   repository.NewFacilityRepo(db),
		Source:     repository.NewEmissionSourceRepo(db),
		Factor:     repository.NewEmissionFactorRepo(db),
		Record:     repository.NewEmissionRecordRepo(db),
		Credit:     repository.NewCarbonCreditRepo(db),
		Target:     repository.NewReductionTargetRepo(db),
		Report:     repository.NewCarbonReportRepo(db),
		Audit:      repository.NewAuditLogRepo(db),
		Analytics:  repository.NewAnalyticsRepo(db),
		User:       repository.NewUserRepo(db),
		DataImport: repository.NewDataImportRepo(db),
		Task:       repository.NewScheduledTaskRepo(db),
		Alert:      repository.NewAlertRepo(db),
		Notify:     repository.NewNotificationRepo(db),
		APIKey:     repository.NewAPIKeyRepo(db),
		Webhook:    repository.NewWebhookRepo(db),
		Attachment: repository.NewAttachmentRepo(db),
		Export:     repository.NewReportExportRepo(db),
		Rollback:   repository.NewRollbackRepo(db),
		Setting:    repository.NewSystemSettingRepo(db),
	}
}

func defaultConfig() *config.Config {
	cfg := &config.Config{}
	cfg.Server.Port = 8080
	cfg.Server.Mode = "debug"
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.SSLMode = "disable"
	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = 6379
	cfg.JWT.Secret = "carbon-emission-management-dev-secret"
	cfg.JWT.ExpireHours = 24
	return cfg
}
