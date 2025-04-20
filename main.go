package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/michielvha/shell-sync/config"
	"github.com/michielvha/shell-sync/filter"
	"github.com/michielvha/shell-sync/syncer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: syncclient <config.yaml>")
		os.Exit(1)
	}
	cfgPath := os.Args[1]
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logging
	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(f)
		defer f.Close()
	}

	// Setup secret filter
	var secretFilter *filter.SecretFilter
	if cfg.Filter.Enabled {
		secretFilter, err = filter.NewSecretFilter(cfg.Filter.Patterns, cfg.Filter.Action)
		if err != nil {
			log.Fatalf("Failed to compile secret filter: %v", err)
		}
	}

	// Handle graceful shutdown
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go syncer.SyncLoop(cfg, secretFilter, stopCh, &wg)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	close(stopCh)
	wg.Wait()
	log.Println("Sync client exited.")
}
