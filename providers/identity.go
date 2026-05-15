package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type IdentityProvider struct {
	Client *http.Client
	APIKey string
}

func NewIdentityProvider() *IdentityProvider {
	// A professional tool NEVER hardcodes API keys.
	apiKey := os.Getenv("IDENTITY_API_KEY")

	return &IdentityProvider{
		Client: &http.Client{
			Timeout: 10 * time.Second, // Hard boundary for network I/O
		},
		APIKey: apiKey,
	}
}

func (i *IdentityProvider) Name() string {
	return "identity_breach_intel"
}

func (i *IdentityProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	if i.APIKey == "" {
		return ProviderResult{}, fmt.Errorf("IDENTITY_API_KEY environment variable is missing")
	}

	// Construct the API endpoint (Replace with your chosen vendor's API)
	endpoint := fmt.Sprintf("https://api.example-breach-db.com/v3/breachedaccount/%s", target)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	// Standard OSINT API headers
	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/2.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.APIKey))
	req.Header.Set("Accept", "application/json")

	resp, err := i.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// 404 in these APIs usually means "clean/not found in breaches"
		return ProviderResult{
			ProviderName: i.Name(),
			Target:       target,
			RawData:      map[string]interface{}{"status": "clean", "breaches": 0},
			RiskScore:    0,
		}, nil
	}

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Parse the JSON response dynamically
	var apiResponse []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Package the raw intel
	rawData := map[string]interface{}{
		"status":   "breached",
		"breaches": len(apiResponse),
		"details":  apiResponse,
	}

	// Base risk calculation logic for identity
	risk := len(apiResponse) * 20
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
