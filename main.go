package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("")
	log.Println("===================================== Netmon =====================================")
	log.Println("Starting netmon...")

	// Parse port flag
	port := flag.Int("port", 9898, "Port to run the web server on")
	flag.Parse()

	// Handle graceful shutdown on interrupt signals
	handleInterrupts()

	// Start data collection in a separate goroutine
	go func() { startCollector() }()

	// Start the web server
	startHttpServer(port)

	log.Println("Exiting...")
	log.Println("===================================================================================")
	log.Println("")
}

func handleInterrupts() {
	log.Println("Press Ctrl+C to stop the server...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()
}
