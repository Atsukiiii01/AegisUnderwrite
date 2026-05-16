package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type IdentityPremiumProvider struct {
	Client *http.Client
	APIKey string
}

func (i *IdentityPremiumProvider) SupportedTypes() []string { return []string{"EMAIL"} }

func NewIdentityPremiumProvider() *IdentityPremiumProvider {
	apiKey := os.Getenv("AEGIS_PREMIUM_IDENTITY_KEY")

	return &IdentityPremiumProvider{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey: apiKey,
	}
}

func (i *IdentityPremiumProvider) Name() string {
	return "identity_premium_intel"
}

func (i *IdentityPremiumProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	if i.APIKey == "" {
		return ProviderResult{}, fmt.Errorf("AEGIS_PREMIUM_IDENTITY_KEY variable is missing")
	}

	// Dynamic boilerplate mapping to a generic premium OSINT proxy/vendor endpoint
	endpoint := fmt.Sprintf("https://api.dehashed.com/v1/search?query=email:%s", target)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/2.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.APIKey))
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
			RawData:      map[string]interface{}{"status": "clean"},
			RiskScore:    0,
		}, nil
	}

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("premium API returned status: %d", resp.StatusCode)
	}

	var parsedResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsedResponse); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to decode structured JSON: %w", err)
	}

	return ProviderResult{
		ProviderName: i.Name(),
		Target:       target,
		RawData:      parsedResponse,
		RiskScore:    95, // Place-holder high severity flag for premium credential detection
	}, nil
}
