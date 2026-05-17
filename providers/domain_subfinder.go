package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DomainSubfinderProvider struct {
	Client *http.Client
}

func NewDomainSubfinderProvider() *DomainSubfinderProvider {
	return &DomainSubfinderProvider{
		Client: &http.Client{
			Timeout: 15 * time.Second, // Passive DNS can take a moment
		},
	}
}

func (p *DomainSubfinderProvider) Name() string {
	return "domain_passive_subfinder"
}

func (p *DomainSubfinderProvider) SupportedTypes() []string {
	return []string{"DOMAIN"}
}

func (p *DomainSubfinderProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	// HackerTarget Free API for passive DNS/Host search
	endpoint := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", target)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/1.1")

	resp, err := p.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	responseStr := string(bodyBytes)

	// HackerTarget returns "error..." on limits or invalid queries
	if strings.HasPrefix(responseStr, "error") {
		return ProviderResult{}, fmt.Errorf("API error/limit reached: %s", responseStr)
	}

	// Parse the CSV output (format: subdomain,ip)
	lines := strings.Split(strings.TrimSpace(responseStr), "\n")
	var subdomains []string

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 0 && parts[0] != "" {
			subdomains = append(subdomains, parts[0])
		}
	}

	subCount := len(subdomains)

	// Risk Logic: A massive footprint implies higher inherent risk of misconfiguration.
	// We cap this module's inherent risk at 35 so it doesn't artificially inflate the score
	// beyond what an actual vulnerability would trigger.
	risk := subCount
	if risk > 35 {
		risk = 35
	}

	// Cap the raw data output to 50 subdomains to avoid blowing up the JSON terminal output
	displaySubs := subdomains
	if len(displaySubs) > 50 {
		displaySubs = displaySubs[:50]
	}

	rawData := map[string]interface{}{
		"total_discovered": subCount,
		"subdomains":       displaySubs,
	}

	return ProviderResult{
		ProviderName: p.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    risk,
	}, nil
}
