package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("🚀 Hytale Status Daemon - Starting...")

	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	log.Printf("Loading config from: %s", *configPath)
	log.Println("Daemon initialized successfully")
}
