package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KelpGF/Go-Expert/08-APIs/configs"
	"github.com/KelpGF/Go-Expert/08-APIs/internal/domain/entity"
	dbRepository "github.com/KelpGF/Go-Expert/08-APIs/internal/infrastructure/database/repository"
	"github.com/KelpGF/Go-Expert/08-APIs/internal/infrastructure/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/KelpGF/Go-Expert/08-APIs/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Go Expert API
// @version         1.0
// @description     Product API with authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Kelvin Gomes
// @contact.url    https://www.linkedin.com/in/kelvin-gomes-fernandes
// @contact.email  kelvingomesdeveloper@gmail.com

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs := configs.LoadConfig(".")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(LogRequest)
	router.Use(middleware.Recoverer)
	router.Use(middleware.WithValue("jwtAuth", configs.TokenAuth))
	router.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go Expert API"))
	})

	router.Get("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("I woke up after 5 seconds"))
	})

	mapperProductRoutes(router, db, configs.TokenAuth)
	mapperUserRoutes(router, db)

	router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/docs/doc.json")))

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

func mapperProductRoutes(router *chi.Mux, db *gorm.DB, jwt *jwtauth.JWTAuth) {
	productRepository := dbRepository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepository)

	router.Route("/product", func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt))
		r.Use(jwtauth.Authenticator)

		r.Get("/{id}", productHandler.Get)
		r.Get("/", productHandler.GetByPagination)
		r.Post("/", productHandler.Create)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)
	})
}

func mapperUserRoutes(router *chi.Mux, db *gorm.DB) {
	userRepository := dbRepository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	router.Post("/user", userHandler.Create)
	router.Post("/user/generate_token", userHandler.GetJwt)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request")
		next.ServeHTTP(w, r)
	})
}
