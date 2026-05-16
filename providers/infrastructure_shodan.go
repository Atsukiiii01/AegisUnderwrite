package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ShodanProvider struct {
	Client *http.Client
	APIKey string
}

func NewShodanProvider() *ShodanProvider {
	apiKey := os.Getenv("SHODAN_API_KEY")

	return &ShodanProvider{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey: apiKey,
	}
}

func (s *ShodanProvider) Name() string {
	return "infrastructure_shodan"
}

// Ensure the router only feeds IP addresses to this module
func (s *ShodanProvider) SupportedTypes() []string {
	return []string{"IP"}
}

func (s *ShodanProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	if s.APIKey == "" {
		return ProviderResult{}, fmt.Errorf("SHODAN_API_KEY environment variable is missing")
	}

	endpoint := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", target, s.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/2.0")
	req.Header.Set("Accept", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return ProviderResult{
			ProviderName: s.Name(),
			Target:       target,
			RawData:      map[string]interface{}{"status": "not_indexed_or_offline"},
			RiskScore:    0,
		}, nil
	}

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("Shodan API returned status: %d", resp.StatusCode)
	}

	var shodanData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&shodanData); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var ports []interface{}
	if p, ok := shodanData["ports"].([]interface{}); ok {
		ports = p
	}

	var vulnerabilities []interface{}
	if vulns, ok := shodanData["vulns"].([]interface{}); ok {
		vulnerabilities = vulns
	}

	orgName := "Unknown"
	if org, ok := shodanData["org"].(string); ok {
		orgName = org
	}

	rawData := map[string]interface{}{
		"organization": orgName,
		"open_ports":   ports,
		"cves":         vulnerabilities,
	}

	// Calculate infrastructure risk:
	// Open ports represent surface area (10 pts each). CVEs represent confirmed holes (50 pts each).
	risk := (len(ports) * 10) + (len(vulnerabilities) * 50)
	if risk > 100 {
		risk = 100
	}

	return ProviderResult{
		ProviderName: s.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    risk,
	}, nil
}
