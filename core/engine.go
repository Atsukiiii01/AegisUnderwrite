package core

import (
	"context"
	"log"
	"sync"

	"github.com/Atsukiiii01/AegisUnderwrite/database"
	"github.com/Atsukiiii01/AegisUnderwrite/providers"
	"github.com/Atsukiiii01/AegisUnderwrite/utils"
)

type Engine struct {
	providers []providers.Provider
	db        *database.Manager // The Engine now owns a database connection
}

// Inject the database manager upon initialization
func NewEngine(dbManager *database.Manager) *Engine {
	return &Engine{
		providers: make([]providers.Provider, 0),
		db:        dbManager,
	}
}

func (e *Engine) Register(p providers.Provider) {
	e.providers = append(e.providers, p)
	log.Printf("[INIT] Registered Provider: %s", p.Name())
}

func (e *Engine) Analyze(ctx context.Context, target string) []providers.ProviderResult {
	targetType := string(utils.IdentifyTarget(target))
	log.Printf("[ROUTER] Identified target '%s' as type: %s\n", target, targetType)

	var wg sync.WaitGroup
	results := make(chan providers.ProviderResult, len(e.providers))

	for _, p := range e.providers {
		supported := false
		for _, t := range p.SupportedTypes() {
			if t == targetType || t == "ANY" {
				supported = true
				break
			}
		}

		if !supported {
			continue
		}

		wg.Add(1)

		go func(provider providers.Provider) {
			defer wg.Done()

			res, err := provider.Fetch(ctx, target)
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", provider.Name(), err)
				return
			}
			results <- res
		}(p)
	}

	wg.Wait()
	close(results)

	var finalReport []providers.ProviderResult

	// Process the results sequentially as they come off the channel
	for r := range results {
		finalReport = append(finalReport, r)

		// Persist to the database
		if e.db != nil {
			err := e.db.SaveAssessment(r.Target, r.ProviderName, r.RawData, r.RiskScore)
			if err != nil {
				log.Printf("[ERROR] Failed to save %s audit to database: %v", r.ProviderName, err)
			} else {
				log.Printf("[DB] Persisted audit record for %s via %s", r.Target, r.ProviderName)
			}
		}
	}

	return finalReport
}
