package presentation

import (
	"log"
	"net/http"
	"time"

	"github.com/KelpGF/Go-Expert/08-APIs/configs"
	dbRepository "github.com/KelpGF/Go-Expert/08-APIs/internal/infrastructure/database/repository"
	"github.com/KelpGF/Go-Expert/08-APIs/internal/infrastructure/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func createRouter(
	db *gorm.DB,
	configs *configs.ConfigType,
) *chi.Mux {

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

	return router
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
