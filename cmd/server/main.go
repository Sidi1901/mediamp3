package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Sidi1901/mediamp3/internal/config"
	"github.com/Sidi1901/mediamp3/internal/router"
)

func main(){

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	cfg, err := config.LoadEnv()
	
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	} 

	addr := ":" + cfg.ServerPort

	r := router.Setup()

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
	
}