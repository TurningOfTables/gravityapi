package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// IsRunninginDocker returns true if the app is running in a docker container, or false if not
// functions by checking for the presence of "./dockerenv"
func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func LoadEnv() {

	if os.Getenv("GRAVITY_API_APP_HOST") != "" {
		log.Println("Env already loaded. Not loading it again!")
		return
	}

	if IsRunningInDocker() {
		godotenv.Load(".env.docker")
		log.Println("Loaded .env.docker")
	} else {
		godotenv.Load(".env.local")
		log.Println("Loaded .env.local")
	}
}
