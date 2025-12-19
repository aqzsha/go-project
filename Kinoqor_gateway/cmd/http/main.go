// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"gateway/api"
	"log"
	"os"
)

func main() {
	log.Println("starting server")

	err := api.NewServer(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
