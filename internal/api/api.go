package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config      config.Config
	router      *gin.Engine
	db          *pgxpool.Pool
	querier     store.Querier
	permCache   *cache.PermissionCache
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	roleHandler *handler.RoleHandler
	jwtSecret   string
}

func NewApp(cfg config.Config, db *pgxpool.Pool) *App {
	querier := store.New(db)
	permCache := cache.NewPermissionCache()

	authSvc := service.NewAuthService(querier, cfg.JWTSecret)
	userSvc := service.NewUserService(querier)
	roleSvc := service.NewRoleService(querier)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)

	app := &App{
		config:      cfg,
		db:          db,
		querier:     querier,
		permCache:   permCache,
		authHandler: authHandler,
		userHandler: userHandler,
		roleHandler: roleHandler,
		jwtSecret:   cfg.JWTSecret,
	}

	app.setupRouter()

	return app
}

func (a *App) setupRouter() {
	router := gin.Default()
	router.Use(cors.Default())
	a.setupRoutes(router)
	a.router = router
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    a.config.Addr,
		Handler: a.router,
	}

	go func() {
		log.Printf("Server starting on %s", a.config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}
