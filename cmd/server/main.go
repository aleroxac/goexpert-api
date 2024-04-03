package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aleroxac/goexpert-api/configs"
	_ "github.com/aleroxac/goexpert-api/docs"
	"github.com/aleroxac/goexpert-api/internal/entity"
	"github.com/aleroxac/goexpert-api/internal/infra/database"
	"github.com/aleroxac/goexpert-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//	@title			goexpert-api
//	@version		1.0.0
//	@description	API criada durante modulo de APIs do treinamento GoExpert da FullCycle
//	@termsOfService	http://swagger.io/terms

//	@contact.name	Augusto Cardoso dos Santos
//	@contact.url	https://github.com/aleroxac
//	@contact.email	acardoso.ti@gmail.com

//	@license.name	Full Cycle License
//	@license.url	https://fullcycle.com.br

// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// configs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	// HTTP Router
	r := chi.NewRouter()

	// middlewares
	// r.Use(LogRequest) 	 	// custom
	r.Use(middleware.Logger)    // go-chi
	r.Use(middleware.Recoverer) // go-chi
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))

	// products
	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)
	r.Route("/products", func(r chi.Router) {
		// middlewares
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)

		// routes
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	// users
	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB)
	r.Route("/users", func(r chi.Router) {
		// routes
		r.Post("/", userHandler.Create)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	// swagger
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	// server
	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
