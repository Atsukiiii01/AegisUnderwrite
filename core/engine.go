package core

import (
	"context"
	"log"
	"sync"

	"github.com/Atsukiiii01/AegisUnderwrite/providers"
	"github.com/Atsukiiii01/AegisUnderwrite/utils"
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
	// 1. Inspect the target archetype
	targetType := string(utils.IdentifyTarget(target))
	log.Printf("[ROUTER] Identified target '%s' as type: %s\n", target, targetType)

	var wg sync.WaitGroup
	results := make(chan providers.ProviderResult, len(e.providers))
	activeWorkers := 0

	for _, p := range e.providers {
		// 2. Check if the provider supports this specific type
		supported := false
		for _, t := range p.SupportedTypes() {
			if t == targetType || t == "ANY" {
				supported = true
				break
			}
		}

		if !supported {
			// Silently skip irrelevant modules without breaking concurrency
			continue
		}

		activeWorkers++
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
	for r := range results {
		finalReport = append(finalReport, r)
	}

	return finalReport
}
