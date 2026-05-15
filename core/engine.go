package core

import (
	"context"
	"log"
	"sync"

	"github.com/Atsukiiii01/AegisUnderwrite/providers"
)

type Engine struct {
	providers []providers.Provider
}

func NewEngine() *Engine {
	return &Engine{
		providers: make([]providers.Provider, 0),
	}
}

func (e *Engine) Register(p providers.Provider) {
	e.providers = append(e.providers, p)
	log.Printf("[INIT] Registered Provider: %s", p.Name())
}

func (e *Engine) Analyze(ctx context.Context, target string) []providers.ProviderResult {
	var wg sync.WaitGroup
	results := make(chan providers.ProviderResult, len(e.providers))

	for _, p := range e.providers {
		wg.Add(1)

		// Spin up a concurrent goroutine for each provider
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

	// Wait for all goroutines to finish, then close the channel
	wg.Wait()
	close(results)

	var finalReport []providers.ProviderResult
	for r := range results {
		finalReport = append(finalReport, r)
	}

	return finalReport
}
