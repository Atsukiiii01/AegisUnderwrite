package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IdentityXONProvider struct {
	Client *http.Client
}

func NewIdentityXONProvider() *IdentityXONProvider {
	return &IdentityXONProvider{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (i *IdentityXONProvider) Name() string {
	return "identity_xposedornot_free"
}

// Ensure the router only feeds emails to this module
func (i *IdentityXONProvider) SupportedTypes() []string {
	return []string{"EMAIL"}
}

func (i *IdentityXONProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	endpoint := fmt.Sprintf("https://api.xposedornot.com/v1/check-email/%s", target)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/2.0")
	req.Header.Set("Accept", "application/json")

	resp, err := i.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return ProviderResult{
			ProviderName: i.Name(),
			Target:       target,
			RawData:      map[string]interface{}{"status": "clean", "breaches": 0},
			RiskScore:    0,
		}, nil
	}

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("API returned unexpected status: %d", resp.StatusCode)
	}

	var breachData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&breachData); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	breachCount := 0
	var sources []string

	if breaches, ok := breachData["breaches"].([]interface{}); ok {
		breachCount = len(breaches)
		for _, b := range breaches {
			if breachDetails, ok := b.([]interface{}); ok && len(breachDetails) > 0 {
				if name, ok := breachDetails[0].(string); ok {
					sources = append(sources, name)
				}
			}
		}
	}

	// --- THE LOGIC FIX IS HERE ---
	// We dynamically assign the status string based on the actual breach count.
	status := "clean"
	if breachCount > 0 {
		status = "breached"
	}

	rawData := map[string]interface{}{
		"status":   status,
		"breaches": breachCount,
		"sources":  sources,
	}

	risk := breachCount * 15
	if risk > 100 {
		risk = 100
	}

	return ProviderResult{
		ProviderName: i.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    risk,
	}, nil
}
