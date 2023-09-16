package main

import (
	"gocker-api/config"
)

func main() {
	server := config.NewAPIServer(":3000")
	server.Run()
}
