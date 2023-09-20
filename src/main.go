package main

import (
	"gocker-api/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal(envErr)
	}

	var listenAddress string
	port, isPresent := os.LookupEnv("PORT")

	if isPresent {
		listenAddress = ":" + port
	} else {
		listenAddress = ":8080"
	}

	server := config.APIServer{ListenAddress: listenAddress}
	log.Printf("Server listening at %s\n", server.ListenAddress)
	log.Fatal(server.Run())
}
