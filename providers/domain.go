package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CrtShEntry maps the JSON structure returned by crt.sh
type CrtShEntry struct {
	NameValue string `json:"name_value"`
}

type DomainProvider struct {
	Client *http.Client
}

func NewDomainProvider() *DomainProvider {
	return &DomainProvider{
		Client: &http.Client{
			Timeout: 15 * time.Second, // crt.sh can be slow, giving it 15 seconds.
		},
	}
}

func (d *DomainProvider) Name() string {
	return "domain_ct_logs"
}

func (d *DomainProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	// Extract the domain if the user accidentally passes a URL or email
	domain := cleanDomainTarget(target)

	endpoint := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/2.0")
	req.Header.Set("Accept", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("crt.sh returned non-200 status: %d", resp.StatusCode)
	}

	var entries []CrtShEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to parse JSON (API might be overloaded): %w", err)
	}

	// Deduplicate subdomains
	uniqueDomains := make(map[string]bool)
	for _, entry := range entries {
		// crt.sh sometimes returns multiple domains separated by newlines
		parts := strings.Split(entry.NameValue, "\n")
		for _, part := range parts {
			cleanPart := strings.TrimSpace(part)
			// Ignore wildcard certs for mapping actual endpoints
			if cleanPart != "" && !strings.HasPrefix(cleanPart, "*.") {
				uniqueDomains[cleanPart] = true
			}
		}
	}

	var subdomains []string
	for dom := range uniqueDomains {
		subdomains = append(subdomains, dom)
	}

	// Risk Calculation: Larger attack surface = higher baseline risk.
	// (1 point per exposed subdomain, capped at 80 for this specific metric)
	riskScore := len(subdomains)
	if riskScore > 80 {
		riskScore = 80
	}

	rawData := map[string]interface{}{
		"base_domain":      domain,
		"subdomains_found": len(subdomains),
		"subdomains":       subdomains,
	}

	return ProviderResult{
		ProviderName: d.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    riskScore,
	}, nil
}

// Utility function to ensure we just query the domain
func cleanDomainTarget(target string) string {
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimPrefix(target, "https://")
	if strings.Contains(target, "@") {
		parts := strings.Split(target, "@")
		return parts[len(parts)-1]
	}
	return target
}
