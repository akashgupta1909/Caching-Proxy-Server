package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func (config *Config) initRouter(redisConfig *RedisConfig) error {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/*", redisConfig.redisClientWrapper(config.handleProxy))

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	fmt.Printf("Starting proxy server of origin %v on port %d\n", config.Origin, config.Port)
	error := server.ListenAndServe()
	if error != nil {
		return fmt.Errorf("error starting server: %v", error)
	}
	return nil
}
