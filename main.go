package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Atsukiiii01/AegisUnderwrite/core"
	"github.com/Atsukiiii01/AegisUnderwrite/providers"
)

func main() {
	log.Println("Starting AegisUnderwrite Engine...")

	engine := core.NewEngine()

	// Register all modules
	engine.Register(providers.NewMockBreachProvider())
	engine.Register(providers.NewIdentityXONProvider())
	engine.Register(providers.NewIdentityPremiumProvider())
	engine.Register(providers.NewDomainProvider())
	engine.Register(providers.NewShodanProvider()) // New Infrastructure Provider

	ctx := context.Background()

	// Test with an IP address to verify routing
	target := "8.8.8.8"

	log.Printf("Executing concurrent analysis against: %s\n", target)
	report := engine.Analyze(ctx, target)

	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("\n--- RESULTS ---")
	os.Stdout.Write(output)
	fmt.Println("\n---------------")
}
