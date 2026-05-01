package main

import (
	"log"
	"os"

	"api-gateway/config"
	"api-gateway/database"
	"api-gateway/model"
	"api-gateway/router"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	seedAdmin(cfg)

	apiRouter := router.SetupAPI(database.DB, cfg)
	proxyRouter := router.SetupProxy(database.DB)

	// Run proxy server in a goroutine
	go func() {
		log.Printf("Proxy server running on :%s", cfg.ProxyPort)
		if err := proxyRouter.Run(":" + cfg.ProxyPort); err != nil {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
	}()

	// Run API server
	log.Printf("API server running on :%s", cfg.ServerPort)
	if err := apiRouter.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}

func seedAdmin(cfg *config.Config) {
	var count int64
	database.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "admin123"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := &model.User{
		Name:     "Admin",
		Email:    "admin@gateway.local",
		Password: string(hashedPassword),
		Role:     "admin",
		IsActive: true,
	}

	if err := database.DB.Create(admin).Error; err != nil {
		log.Printf("Failed to seed admin user: %v", err)
		return
	}

	log.Println("Admin user created: admin@gateway.local / admin123")
}
