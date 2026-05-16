package core

import (
	"context"
	"log"
	"sync"

	"github.com/Atsukiiii01/AegisUnderwrite/database"
	"github.com/Atsukiiii01/AegisUnderwrite/providers"
	"github.com/Atsukiiii01/AegisUnderwrite/utils"
)

type AssessmentReport struct {
	Target      string                     `json:"target"`
	TargetType  string                     `json:"target_type"`
	OverallRisk int                        `json:"overall_risk"`
	RiskLevel   string                     `json:"risk_level"`
	Findings    []providers.ProviderResult `json:"findings"`
}

type Engine struct {
	providers []providers.Provider
	db        *database.Manager
}

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

func (e *Engine) Analyze(ctx context.Context, target string) AssessmentReport {
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

	var findings []providers.ProviderResult
	maxScore := 0
	totalScore := 0
	validProviders := 0

	for r := range results {
		findings = append(findings, r)

		if r.RiskScore > maxScore {
			maxScore = r.RiskScore
		}
		totalScore += r.RiskScore
		validProviders++

		if e.db != nil {
			err := e.db.SaveAssessment(r.Target, r.ProviderName, r.RawData, r.RiskScore)
			if err != nil {
				log.Printf("[ERROR] Failed to save %s audit to database: %v", r.ProviderName, err)
			}
		}
	}

	// The Underwriter Algorithm:
	// Base risk equals the absolute highest threat found.
	// Compounding factor: add 10% of the remaining combined scores.
	finalRisk := maxScore
	if validProviders > 1 {
		remainder := totalScore - maxScore
		finalRisk += int(float64(remainder) * 0.1)
	}
	if finalRisk > 100 {
		finalRisk = 100
	}

	level := "LOW"
	if finalRisk >= 30 {
		level = "MEDIUM"
	}
	if finalRisk >= 70 {
		level = "HIGH"
	}
	if finalRisk >= 90 {
		level = "CRITICAL"
	}

	return AssessmentReport{
		Target:      target,
		TargetType:  targetType,
		OverallRisk: finalRisk,
		RiskLevel:   level,
		Findings:    findings,
	}
}
