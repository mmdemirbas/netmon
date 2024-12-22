package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

var (
	absoluteDbFilePath string         // Absolute path to the database file
	collectorInterval  *time.Duration // Interval between measurements
	serverPort         *int           // Port to run the web server on
	svcFlag            *string        // Control the system service: install, uninstall, start, stop, restart

	logger service.Logger

	// Configuration
	serviceName        = "netmon"
	serviceDescription = "Network Monitoring Service - github.com/mmdemirbas/netmon"
)

// TODO: Graceful shutdown collector

func main() {
	// Parse command-line flags
	dbFileName := flag.String("db-file", "data/netmon.db", "(optional) Database file name.")
	collectorInterval = flag.Duration("interval", 5*60*1_000_000_000, "(optional) Interval between measurements.")
	serverPort = flag.Int("port", 9898, "(optional) Port to run the web server on.")
	svcFlag = flag.String("service", "", "(optional) Control the system service: install, uninstall, start, stop, restart")
	flag.Parse()

	var err error
	absoluteDbFilePath, err = filepath.Abs(*dbFileName)

	// Initialize the service infra
	svc := initProgram()

	// Handle service control actions (service install, uninstall, start, stop)
	if len(*svcFlag) != 0 {
		// TODO: Honor cli flags during service installation or start maybe
		err = service.Control(svc, *svcFlag)
		if err == nil {
			logger.Infof("Service %s succeeded", *svcFlag)
			os.Exit(0)
		} else if !slices.Contains(service.ControlAction[:], *svcFlag) {
			logger.Errorf("Service %s failed: Valid actions are: %v", *svcFlag, service.ControlAction)
			os.Exit(1)
		} else {
			logger.Errorf("Service %s failed: %v", *svcFlag, err)
			os.Exit(1)
		}
	}

	// Run the service
	err = svc.Run()
	if err != nil {
		logger.Errorf("Failed to run service: %v", err)
	}
}

type program struct {
	exit chan struct{}
}

func initProgram() service.Service {
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}

	// Create a new program instance
	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Error creating service: %v\n", err)
	}

	// Initialize the logger
	errs := make(chan error, 5)
	logger, err = svc.Logger(errs)
	if err != nil {
		log.Fatalf("Error getting logger: %v\n", err)
	}

	// Handle errors from the logger
	go func() {
		for {
			err = <-errs
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	return svc
}

func (p *program) Start(_ service.Service) error {
	if service.Interactive() {
		// Running in terminal
	} else {
		// Running as a service
	}

	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go func() {
		err := p.run()
		if err != nil {
			fmt.Errorf("Error running service: %v", err)
			close(p.exit)
		}
	}()

	return nil
}

// actual logic of the service
func (p *program) run() error {
	logger.Infof("===================================== Netmon =====================================")
	logger.Infof("Starting netmon with the following settings:")
	logger.Infof("  -db-file   = %s", absoluteDbFilePath)
	logger.Infof("  -interval  = %v", *collectorInterval)
	logger.Infof("  -port      = %d", *serverPort)

	// Initialize the database
	err := initDatabase(absoluteDbFilePath)
	if err != nil {
		return err
	}

	// Start data collection in a separate goroutine
	ticker := startCollector(collectorInterval)
	defer ticker.Stop()

	// Start the web server
	return startHttpServer(serverPort)
}

func (p *program) Stop(_ service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}
