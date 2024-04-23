package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pttrulez/investor-go/internal/config"
	"github.com/pttrulez/investor-go/internal/controllers"
	"github.com/pttrulez/investor-go/internal/repository/postgres"
	"github.com/pttrulez/investor-go/internal/services"
)

func main() {
	cfg := config.MustLoad()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	repository, err := postgres.NewPostgresRepo(cfg.Pg)
	if err != nil {
		panic("Failed to initialize postgres repository: " + err.Error())
	}
	
	services := services.NewServiceContainer(repository)
	controllers.Init(r, repository, services)

	logger := slog.Default()
	address := fmt.Sprintf("%v:%v", cfg.ApiHost, cfg.ApiPort)
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}
	
	logger.Info(fmt.Sprintf("Listening on  %v", address))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}
