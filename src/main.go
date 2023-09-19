package main

import (
	"gocker-api/config"
	"log"
)

func main() {
	server := config.APIServer{ListenAddress: ":8080"}
	log.Printf("Server listening at %s\n", server.ListenAddress)
	log.Fatal(server.Run())
}
