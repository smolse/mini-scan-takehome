package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/smolse/scan-takehome/internal/config"
	"github.com/smolse/scan-takehome/internal/datastores"
	"github.com/smolse/scan-takehome/internal/services"
)

func main() {
	// Load the application configuration
	cfg, err := config.LoadProcessorConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and connect to the data store
	db, err := datastores.NewScanDataStore(&cfg.DataStore)
	if err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}
	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to data store: %v", err)
	}
	defer db.Close()

	// Initialize the processor service, passing the data store as a dependency
	svc := services.NewProcessorService(db)

	// Initialize the Pub/Sub client
	pubsubClient, err := pubsub.NewClient(context.Background(), cfg.PubSub.ProjectId)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub client: %v", err)
	}
	defer pubsubClient.Close()

	// Initialize and configure the Pub/Sub subscription
	subscription := pubsubClient.Subscription(cfg.PubSub.SubscriptionId)
	subscription.ReceiveSettings.MaxOutstandingMessages = cfg.PubSub.MaxOutstandingMessages

	// Gracefully handle SIGINT and SIGTERM signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a context to cancel the Pub/Sub message receiving loop
	ctx, cancel := context.WithCancel(context.Background())

	// Start polling and processing messages from Pub/Sub
	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			err := svc.ProcessScanData(msg.Data)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		})
		if err != nil {
			log.Fatalf("Failed to receive messages: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down service...")

	// Cancel the context to stop receiving new messages
	cancel()

	// Allow time for outstanding messages to be acknowledged
	time.Sleep(cfg.Service.GracefulShutdownTimeout)
	log.Println("Shutdown completed!")
}
