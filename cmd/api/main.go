package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fesbarbosa/melivendas-api/internal/adapters/input/http/handlers"
	"github.com/fesbarbosa/melivendas-api/internal/adapters/input/http/routes"
	"github.com/fesbarbosa/melivendas-api/internal/adapters/output/db"
	"github.com/fesbarbosa/melivendas-api/internal/config"
	"github.com/fesbarbosa/melivendas-api/internal/core/services"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.NewConfig()

	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Falha ao inicializar banco de dados: %v", err)
	}
	defer database.Close()

	itemRepository := db.NewItemRepository(database)

	itemService := services.NewItemService(itemRepository)

	itemHandler := handlers.NewItemHandler(itemService)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	routes.RegisterItemRoutes(router, itemHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Servidor escutando na porta %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Falha ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Servidor forçado a desligar: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
