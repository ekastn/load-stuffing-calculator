package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"load-stuffing-calculator/internal/config"
	"load-stuffing-calculator/internal/gateway"
	"load-stuffing-calculator/internal/handler"
	"load-stuffing-calculator/internal/service"
	"load-stuffing-calculator/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App menampung semua dependencies aplikasi.
// Dengan menyimpan semua komponen di satu struct, kita bisa dengan mudah
// melihat apa saja yang dibutuhkan aplikasi untuk berjalan.
type App struct {
	config           *config.Config
	router           *gin.Engine
	db               *pgxpool.Pool
	containerHandler *handler.ContainerHandler
	productHandler   *handler.ProductHandler
	planHandler      *handler.PlanHandler
}

// NewApp menginisialisasi semua layers dan dependencies.
// Urutan inisialisasi mengikuti dependency graph:
// database → store → services → handlers → router
func NewApp(cfg *config.Config, db *pgxpool.Pool) *App {
	// Store layer: wrapper untuk database queries
	querier := store.New(db)

	// Initialize gateway untuk Packing Service
	// Gunakan MockPackingGateway untuk demo tanpa Packing Service,
	// atau HTTPPackingGateway untuk koneksi ke service yang sebenarnya
	// packingGW := gateway.NewMockPackingGateway()

	// Untuk production, gunakan:
	packingGW := gateway.NewHTTPPackingGateway(
		cfg.PackingServiceURL,
		60*time.Second,
	)

	// Service layer: business logic
	containerSvc := service.NewContainerService(querier)
	productSvc := service.NewProductService(querier)
	planSvc := service.NewPlanService(querier, packingGW)

	// Handler layer: HTTP request/response handling
	containerHandler := handler.NewContainerHandler(containerSvc)
	productHandler := handler.NewProductHandler(productSvc)
	planHandler := handler.NewPlanHandler(planSvc)

	app := &App{
		config:           cfg,
		db:               db,
		containerHandler: containerHandler,
		productHandler:   productHandler,
		planHandler:      planHandler,
	}

	app.setupRouter()

	return app
}

func (a *App) setupRouter() {
	// gin.Default() membuat router dengan Logger dan Recovery middleware
	router := gin.Default()

	// CORS middleware memungkinkan frontend dari domain berbeda
	// untuk mengakses API kita
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Setup routes di file terpisah untuk organisasi yang lebih baik
	a.setupRoutes(router)
	a.router = router
}

func (a *App) Run() error {
	// Buat http.Server dengan konfigurasi kustom
	srv := &http.Server{
		Addr:    ":" + a.config.ServerPort,
		Handler: a.router,
	}

	// Jalankan server di goroutine terpisah agar tidak blocking
	go func() {
		log.Printf("Server starting on :%s", a.config.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown: tunggu signal interrupt (Ctrl+C) atau terminate
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block sampai menerima signal
	log.Println("Shutting down server...")

	// Beri waktu 10 detik untuk request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}
