package main

import (
	"github.com/KelpGF/Go-Expert/08-APIs/configs"
	"github.com/KelpGF/Go-Expert/08-APIs/internal/domain/entity"
	"github.com/KelpGF/Go-Expert/08-APIs/internal/presentation"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/KelpGF/Go-Expert/08-APIs/docs"
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

	presentation.StartWebServer(db, configs)
}
