package main

import (
	"log"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/handler"
)

func main() {

	c, err := config.New()
	if err != nil {
		log.Fatal("failed to initialize configuration: ", err)
	}

	h := handler.New(c)
	if err := h.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
