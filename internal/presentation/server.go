package presentation

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KelpGF/Go-Expert/08-APIs/configs"
	"gorm.io/gorm"
)

func StartWebServer(
	db *gorm.DB,
	configs *configs.ConfigType,
) {
	router := createRouter(db, configs)
	server := &http.Server{
		Addr:    configs.WebServerHost + ":" + configs.WebServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server is running on http://%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Server is shutting down...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	log.Println("Server stopped")
}
