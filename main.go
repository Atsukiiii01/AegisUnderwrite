package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Atsukiiii01/AegisUnderwrite/core"
	"github.com/Atsukiiii01/AegisUnderwrite/database"
	"github.com/Atsukiiii01/AegisUnderwrite/providers"
)

func main() {
	// 1. Define Subcommands
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	scanTarget := scanCmd.String("target", "", "Target to analyze (Email, IP, Domain)")

	historyCmd := flag.NewFlagSet("history", flag.ExitOnError)
	historyTarget := historyCmd.String("target", "", "Target to retrieve history for")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'scan' or 'history' subcommands.")
		fmt.Println("Usage: ./aegis scan -target <target>")
		os.Exit(1)
	}

	// 2. Initialize Database
	dbManager, err := database.NewManager("aegis_audit.db")
	if err != nil {
		log.Fatalf("FATAL: Could not initialize database: %v", err)
	}

	// 3. Route Subcommands
	switch os.Args[1] {
	case "scan":
		scanCmd.Parse(os.Args[2:])
		if *scanTarget == "" {
			fmt.Println("Error: -target is required for scan")
			scanCmd.PrintDefaults()
			os.Exit(1)
		}
		executeScan(dbManager, *scanTarget)

	case "history":
		historyCmd.Parse(os.Args[2:])
		if *historyTarget == "" {
			fmt.Println("Error: -target is required for history")
			historyCmd.PrintDefaults()
			os.Exit(1)
		}
		executeHistory(dbManager, *historyTarget)

	default:
		fmt.Println("Expected 'scan' or 'history' subcommands")
		os.Exit(1)
	}
}

func executeScan(db *database.Manager, target string) {
	log.Printf("[SYSTEM] Booting AegisUnderwrite Engine against: %s", target)
	engine := core.NewEngine(db)

	engine.Register(providers.NewMockBreachProvider())
	engine.Register(providers.NewIdentityXONProvider())
	engine.Register(providers.NewIdentityPremiumProvider())
	engine.Register(providers.NewDomainProvider())
	engine.Register(providers.NewShodanProvider())

	report := engine.Analyze(context.Background(), target)

	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("\n--- FINAL REPORT ---")
	os.Stdout.Write(output)
	fmt.Println("\n--------------------")
}

func executeHistory(db *database.Manager, target string) {
	log.Printf("[SYSTEM] Retrieving audit history for: %s", target)
	records, err := db.GetHistory(target)
	if err != nil {
		log.Fatalf("Failed to fetch history: %v", err)
	}

	if len(records) == 0 {
		fmt.Printf("No history found for target: %s\n", target)
		return
	}

	for _, r := range records {
		fmt.Printf("[%s] %s | Provider: %s | Risk Score: %d\n", r.Timestamp, r.Target, r.Provider, r.RiskScore)
	}
}
