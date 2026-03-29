package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"clash-server/internal/config"
	"clash-server/internal/handler"
	"clash-server/internal/middleware"
	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/internal/scheduler"
	"clash-server/internal/service"
	"clash-server/res"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {
	serverCfg := config.InitServerConfig()
	if err := model.InitDB(serverCfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	coreRepo := repository.NewConfigRepository()
	config.InitCoreRepository(coreRepo)
	hub := handler.InitWebSocketHub()

	coreService := service.GetCoreService()
	coreService.SetHub(hub)
	coreService.SetOnStatusChange(func(status service.CoreStatus) {
		handler.BroadcastCoreStatus(status.Running, status.Version, status.Error)
	})
	if err := coreService.Start(); err != nil {
		log.Printf("Warning: Failed to start core: %v", err)
	}

	subScheduler := scheduler.NewSubscriptionScheduler()
	if err := subScheduler.Start(); err != nil {
		log.Printf("Warning: Failed to start subscription scheduler: %v", err)
	}
	handler.SetSubscriptionScheduler(subScheduler)

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(sessions.Sessions(
		"session",
		memstore.NewStore([]byte(generateSessionSecret())),
	))
	api := r.Group("/api")
	{
		authHandler := handler.NewAuthHandler()
		api.GET("/init", authHandler.CheckInit)
		api.POST("/init", authHandler.InitPassword)
		api.POST("/session", authHandler.Login)
		api.DELETE("/session", authHandler.Logout)
		api.GET("/version", handler.GetVersion)
	}
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		userHandler := handler.NewUserHandler()
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me/password", userHandler.ChangePassword)
		subHandler := handler.NewSubscriptionHandler()
		protected.GET("/subscriptions", subHandler.List)
		protected.POST("/subscriptions", subHandler.Create)
		protected.GET("/subscriptions/:id", subHandler.Get)
		protected.PUT("/subscriptions/:id", subHandler.Update)
		protected.DELETE("/subscriptions/:id", subHandler.Delete)
		protected.POST("/subscriptions/:id/refresh", subHandler.Refresh)
		protected.PUT("/subscriptions/:id/activate", subHandler.Activate)
		protected.GET("/subscriptions/:id/merged", subHandler.GetMerged)
		protected.GET("/subscriptions/:id/content", subHandler.GetContent)
		protected.PUT("/subscriptions/:id/content", subHandler.UpdateContent)
		customizationHandler := handler.NewCustomizationHandler()
		protected.GET("/subscriptions/:id/customization", customizationHandler.Get)
		protected.PUT("/subscriptions/:id/customization", customizationHandler.Update)
		configHandler := handler.NewConfigHandler()
		protected.GET("/config", configHandler.Get)
		protected.PUT("/config", configHandler.Update)
		systemHandler := handler.NewSystemHandler()
		protected.GET("/system/info", systemHandler.GetInfo)
		proxyHandler := handler.NewProxyHandler()
		protected.GET("/proxies", proxyHandler.List)
		protected.GET("/proxies/:name", proxyHandler.Get)
		protected.PUT("/proxies/:name", proxyHandler.Select)
		protected.GET("/proxies/:name/delay", proxyHandler.CheckDelay)
		protected.GET("/proxies/group/:group/delay", proxyHandler.CheckGroupDelay)
		protected.GET("/proxies/mode", proxyHandler.GetMode)
		protected.PATCH("/proxies/mode", proxyHandler.SetMode)
		connHandler := handler.NewConnectionHandler()
		protected.GET("/connections", connHandler.List)
		protected.DELETE("/connections", connHandler.CloseAll)
		protected.DELETE("/connections/:id", connHandler.Close)
		ruleHandler := handler.NewRuleHandler()
		protected.GET("/rules", ruleHandler.GetRules)
		protected.GET("/ws", handler.WebSocketHandler)
	}

	faviconRegex := regexp.MustCompile("/favicon.*")
	assetsRegex := regexp.MustCompile("/assets/.*")
	r.NoRoute(func(ctx *gin.Context) {
		name := ctx.Request.URL.Path
		if faviconRegex.MatchString(ctx.Request.URL.Path) {
			name = name[1:]
			http.ServeFileFS(ctx.Writer, ctx.Request, res.WebFS, name)
			return
		}
		if assetsRegex.MatchString(ctx.Request.URL.Path) {
			name = name[8:]
			http.ServeFileFS(ctx.Writer, ctx.Request, res.WebAssetsFS, name)
			return
		}
		http.ServeFileFS(ctx.Writer, ctx.Request, res.WebFS, "index.html")
	})

	go func() {
		addr := serverCfg.Host + ":" + strconv.Itoa(serverCfg.Port)
		log.Printf("Server starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	subScheduler.Stop()
	coreService.Stop()
	log.Println("Server stopped")
}

func generateSessionSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate session secret: %v", err)
	}
	secret := hex.EncodeToString(bytes)
	log.Printf("Generated session secret: %s", secret)
	return secret
}
