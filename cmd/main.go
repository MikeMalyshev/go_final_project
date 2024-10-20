package main

import (
	"log"
	"net/http"

	"go_final_project/internal/app"
	"go_final_project/internal/config"
	"go_final_project/internal/db"
	"go_final_project/internal/rest"
)

func main() {
	log.Println("Starting server")

	cfg := config.New()
	database := db.New(cfg)
	if !database.Exists() {
		if err := database.Create(); err != nil {
			log.Fatal("main(): %w ", err)
			return
		}
	}

	if err := database.Open(); err != nil {
		log.Fatal("main(): %w ", err)
		return
	}

	application := app.CreateApplication(database)
	mux := rest.NewMux(application, cfg)

	log.Printf("Server is running on port: :%s\n", cfg.Port())

	err := http.ListenAndServe(":"+cfg.Port(), mux.ServeMux())
	if err != nil {
		log.Fatal("main(): %w ", err)
	}
}
