package router

import (
	"api-gateway/config"
	"api-gateway/handler"
	"api-gateway/middleware"
	"api-gateway/proxy"
	"api-gateway/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAPI(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	userRepo := repository.NewUserRepository(db)
	serviceRepo := repository.NewServiceRepository(db)

	authHandler := handler.NewAuthHandler(userRepo, cfg)
	userHandler := handler.NewUserHandler(userRepo)
	serviceHandler := handler.NewServiceHandler(serviceRepo)

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
		api.GET("/gateway/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success":        true,
				"proxy_base_url": cfg.ProxyBaseURL,
			})
		})
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/auth/me", authHandler.Me)

		// User management (admin only)
		users := protected.Group("/users")
		users.Use(middleware.AdminOnly())
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.PUT("/:id/reset-password", userHandler.ResetPassword)
		}

		// Service management (admin only)
		services := protected.Group("/services")
		services.Use(middleware.AdminOnly())
		{
			services.POST("", serviceHandler.Create)
			services.GET("", serviceHandler.GetAll)
			services.GET("/:id", serviceHandler.GetByID)
			services.PUT("/:id", serviceHandler.Update)
			services.DELETE("/:id", serviceHandler.Delete)
		}
	}

	return r
}

func SetupProxy(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	serviceRepo := repository.NewServiceRepository(db)
	proxyHandler := proxy.NewProxyHandler(serviceRepo)

	// All requests go through the proxy
	r.NoRoute(proxyHandler.Handle)

	return r
}
