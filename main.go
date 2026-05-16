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

	// Register modules explicitly
	engine.Register(providers.NewMockBreachProvider())
	engine.Register(providers.NewIdentityXONProvider())     // Free
	engine.Register(providers.NewIdentityPremiumProvider()) // Premium
	engine.Register(providers.NewDomainProvider())

	ctx := context.Background()
	target := "test@gmail.com"

	log.Printf("Executing concurrent analysis against: %s\n", target)
	report := engine.Analyze(ctx, target)

	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("\n--- RESULTS ---")
	os.Stdout.Write(output)
	fmt.Println("\n---------------")
}
