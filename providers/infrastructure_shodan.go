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
	return &ShodanProvider{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey: os.Getenv("SHODAN_API_KEY"),
	}
}

func (s *ShodanProvider) Name() string {
	return "infrastructure_shodan"
}

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

	req.Header.Set("User-Agent", "AegisUnderwrite-Engine/1.2")

	resp, err := s.Client.Do(req)
	if err != nil {
		return ProviderResult{}, fmt.Errorf("network execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// 404 just means Shodan hasn't scanned this IP. It is not an error, it's a clean result.
		return ProviderResult{
			ProviderName: s.Name(),
			Target:       target,
			RawData:      map[string]interface{}{"status": "not_in_database", "open_ports": []int{}, "cves": []string{}},
			RiskScore:    0,
		}, nil
	}

	if resp.StatusCode != 200 {
		return ProviderResult{}, fmt.Errorf("API returned unexpected status: %d", resp.StatusCode)
	}

	// Unmarshal only the fields we care about for the underwriter model
	var data struct {
		Org   string   `json:"org"`
		Ports []int    `json:"ports"`
		Vulns []string `json:"vulns"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ProviderResult{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// --- UPGRADED RISK LOGIC ---
	// Base penalty: 5 points per open port.
	// Critical penalty: 35 points per known CVE.
	risk := len(data.Ports) * 5

	if len(data.Vulns) > 0 {
		risk += len(data.Vulns) * 35
	}

	if risk > 100 {
		risk = 100
	}

	// Ensure we don't output 'null' for empty CVE slices in the final JSON
	cves := data.Vulns
	if cves == nil {
		cves = []string{}
	}

	rawData := map[string]interface{}{
		"organization": data.Org,
		"open_ports":   data.Ports,
		"cves":         cves,
	}

	return ProviderResult{
		ProviderName: s.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    risk,
	}, nil
}
