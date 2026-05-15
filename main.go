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

	// Register modules
	engine.Register(providers.NewMockBreachProvider())
	engine.Register(providers.NewIdentityProvider()) // The new live provider

	ctx := context.Background()
	target := "target@example.com" // Provide a real email to test

	log.Printf("Executing concurrent analysis against: %s\n", target)
	report := engine.Analyze(ctx, target)

	// Output clean JSON
	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("\n--- RESULTS ---")
	os.Stdout.Write(output)
	fmt.Println("\n---------------")
}
