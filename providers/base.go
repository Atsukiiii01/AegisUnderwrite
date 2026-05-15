package providers

import "context"

// Provider defines the strict interface all OSINT sources must follow.
type Provider interface {
	Name() string
	Fetch(ctx context.Context, target string) (ProviderResult, error)
}

// ProviderResult standardizes the output across all disparate data sources.
type ProviderResult struct {
	ProviderName string                 `json:"provider_name"`
	Target       string                 `json:"target"`
	RawData      map[string]interface{} `json:"raw_data"`
	RiskScore    int                    `json:"risk_score"`
}