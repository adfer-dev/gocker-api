package main

import (
	"gocker-api/config"
)

func main() {
	server := config.NewAPIServer(":8080")
	server.Run()
}
