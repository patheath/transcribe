package main

import (
	"context"
	"log"
	"os"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/cloudevents/sdk-go/v2/event"
	_ "github.com/patheath/transcribe"
)

func main() {

	if err := funcframework.RegisterCloudEventFunctionContext(context.Background(), "/test", func(ctx context.Context, e event.Event) error {
		log.Printf("Received CloudEvent: %v", e)
		return nil
	}); err != nil {
		log.Fatalf("funcframework.RegisterCloudEventFunctionContext: %v\n", err)
	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with 
	// LOCAL_ONLY=true to avoid triggering firewall warnings and 
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	} 
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}