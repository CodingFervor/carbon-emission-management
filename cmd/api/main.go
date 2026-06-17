package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/cache"
	"github.com/CodingFervor/carbon-emission-management/internal/config"
	"github.com/CodingFervor/carbon-emission-management/internal/database"
	"github.com/CodingFervor/carbon-emission-management/internal/middleware"
	"github.com/CodingFervor/carbon-emission-management/pkg/jwt"
	"github.com/CodingFervor/carbon-emission-management/pkg/logger"
)

func main() {
	// Load configuration (falls back to sensible defaults if file is absent).
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Warn("config file not loaded, using defaults", "error", err)
		cfg = &config.Config{}
		cfg.Server.Port = 8080
		cfg.Server.Mode = "debug"
		cfg.Database.Host = "localhost"
		cfg.Database.Port = 5432
		cfg.Database.SSLMode = "disable"
		cfg.Redis.Host = "localhost"
		cfg.Redis.Port = 6379
		cfg.JWT.Secret = "carbon-emission-management-dev-secret"
		cfg.JWT.ExpireHours = 24
	}

	gin.SetMode(cfg.Server.Mode)
	logger.SetLevel(cfg.Server.Mode)

	// Wire up the JWT secret for token signing/verification.
	jwt.SetSecret(cfg.JWT.Secret)

	// Connect infrastructure dependencies. Failures are logged but do not
	// abort startup so the API can still serve health/liveness probes.
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

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")
	{
		// Auth (public)
		api.POST("/auth/login", Login)
		api.POST("/auth/register", Register)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			// Auth (authenticated)
			auth.GET("/auth/profile", Profile)

			// Organizations
			auth.GET("/organizations", ListOrganizations)
			auth.POST("/organizations", CreateOrganization)
			auth.GET("/organizations/:id", GetOrganization)
			auth.PUT("/organizations/:id", UpdateOrganization)
			auth.DELETE("/organizations/:id", DeleteOrganization)

			// Facilities
			auth.GET("/facilities", ListFacilities)
			auth.POST("/facilities", CreateFacility)
			auth.GET("/facilities/:id", GetFacility)
			auth.PUT("/facilities/:id", UpdateFacility)
			auth.DELETE("/facilities/:id", DeleteFacility)
			auth.GET("/organizations/:id/facilities", ListFacilitiesByOrg)

			// Emission Sources
			auth.GET("/emission-sources", ListEmissionSources)
			auth.POST("/emission-sources", CreateEmissionSource)
			auth.GET("/emission-sources/:id", GetEmissionSource)
			auth.PUT("/emission-sources/:id", UpdateEmissionSource)
			auth.DELETE("/emission-sources/:id", DeleteEmissionSource)

			// Emission Factors
			auth.GET("/emission-factors", ListEmissionFactors)
			auth.POST("/emission-factors", CreateEmissionFactor)
			auth.GET("/emission-factors/:id", GetEmissionFactor)
			auth.PUT("/emission-factors/:id", UpdateEmissionFactor)
			auth.DELETE("/emission-factors/:id", DeleteEmissionFactor)

			// Emission Records
			auth.GET("/emission-records", ListEmissionRecords)
			auth.POST("/emission-records", CreateEmissionRecord)
			auth.GET("/emission-records/:id", GetEmissionRecord)
			auth.PUT("/emission-records/:id", UpdateEmissionRecord)
			auth.DELETE("/emission-records/:id", DeleteEmissionRecord)
			auth.GET("/facilities/:id/emission-records", ListRecordsByFacility)
			auth.POST("/emission-records/calculate", CalculateEmission)

			// Carbon Credits
			auth.GET("/carbon-credits", ListCarbonCredits)
			auth.POST("/carbon-credits", CreateCarbonCredit)
			auth.GET("/carbon-credits/:id", GetCarbonCredit)
			auth.PUT("/carbon-credits/:id", UpdateCarbonCredit)
			auth.DELETE("/carbon-credits/:id", DeleteCarbonCredit)
			auth.POST("/carbon-credits/:id/retire", RetireCarbonCredit)

			// Reduction Targets
			auth.GET("/reduction-targets", ListReductionTargets)
			auth.POST("/reduction-targets", CreateReductionTarget)
			auth.GET("/reduction-targets/:id", GetReductionTarget)
			auth.PUT("/reduction-targets/:id", UpdateReductionTarget)
			auth.DELETE("/reduction-targets/:id", DeleteReductionTarget)

			// Carbon Reports
			auth.GET("/reports", ListReports)
			auth.POST("/reports", CreateReport)
			auth.GET("/reports/:id", GetReport)
			auth.PUT("/reports/:id", UpdateReport)
			auth.DELETE("/reports/:id", DeleteReport)
			auth.POST("/reports/:id/generate", GenerateReport)

			// Analytics
			auth.GET("/analytics/dashboard", Dashboard)
			auth.GET("/analytics/by-scope", EmissionsByScope)
			auth.GET("/analytics/trend", EmissionsTrend)
			auth.GET("/analytics/comparison", EmissionsComparison)
			auth.GET("/analytics/facility-breakdown", FacilityBreakdown)

			// Audit Logs (admin only)
			auth.GET("/audit-logs", middleware.RequireRole("admin"), ListAuditLogs)
		}
	}

	addr := ":" + itoa(cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

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

// itoa is a dependency-free int->string helper to keep the binary light.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [12]byte{}
	pos := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// ---------------------------------------------------------------------------
// Handler stubs. Each returns a canned response matching the project's
// response conventions; wire them to real persistence as needed.
// ---------------------------------------------------------------------------

// --- Auth ---

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.example.token",
		"user":  gin.H{"id": 1, "username": "admin", "role": "admin"},
	})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id": 1, "username": "admin", "email": "admin@carbon.io", "role": "admin",
	}})
}

// --- Organizations ---

func ListOrganizations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateOrganization(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "organization created"})
}
func GetOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "organization updated"})
}
func DeleteOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "organization deleted"})
}

// --- Facilities ---

func ListFacilities(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateFacility(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "facility created"})
}
func GetFacility(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateFacility(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "facility updated"})
}
func DeleteFacility(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "facility deleted"})
}
func ListFacilitiesByOrg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}

// --- Emission Sources ---

func ListEmissionSources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateEmissionSource(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "emission source created"})
}
func GetEmissionSource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateEmissionSource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission source updated"})
}
func DeleteEmissionSource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission source deleted"})
}

// --- Emission Factors ---

func ListEmissionFactors(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateEmissionFactor(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "emission factor created"})
}
func GetEmissionFactor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateEmissionFactor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission factor updated"})
}
func DeleteEmissionFactor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission factor deleted"})
}

// --- Emission Records ---

func ListEmissionRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateEmissionRecord(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "emission record created"})
}
func GetEmissionRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateEmissionRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission record updated"})
}
func DeleteEmissionRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "emission record deleted"})
}
func ListRecordsByFacility(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CalculateEmission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"co2_kg": 0, "calculated": true}})
}

// --- Carbon Credits ---

func ListCarbonCredits(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateCarbonCredit(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "carbon credit created"})
}
func GetCarbonCredit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateCarbonCredit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "carbon credit updated"})
}
func DeleteCarbonCredit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "carbon credit deleted"})
}
func RetireCarbonCredit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "carbon credit retired"})
}

// --- Reduction Targets ---

func ListReductionTargets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateReductionTarget(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "reduction target created"})
}
func GetReductionTarget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateReductionTarget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "reduction target updated"})
}
func DeleteReductionTarget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "reduction target deleted"})
}

// --- Carbon Reports ---

func ListReports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func CreateReport(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "report created"})
}
func GetReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
func UpdateReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "report updated"})
}
func DeleteReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "report deleted"})
}
func GenerateReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "report generated",
		"data":    gin.H{"total_co2_t": 0, "scope1_co2_t": 0, "scope2_co2_t": 0, "scope3_co2_t": 0, "status": "generated"},
	})
}

// --- Analytics ---

func Dashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"total_co2_t":     0,
		"scope1_co2_t":    0,
		"scope2_co2_t":    0,
		"scope3_co2_t":    0,
		"offsets_t":       0,
		"net_co2_t":       0,
		"yoy_change_pct":  0,
	}})
}
func EmissionsByScope(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"scope_1": 0, "scope_2": 0, "scope_3": 0,
	}})
}
func EmissionsTrend(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
func EmissionsComparison(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"baseline": 0, "current": 0, "reduction_pct": 0}})
}
func FacilityBreakdown(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}

// --- Audit Logs ---

func ListAuditLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
}
