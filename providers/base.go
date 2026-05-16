package providers

import "context"

// Provider defines the strict interface all OSINT sources must follow.
type Provider interface {
	Name() string
	// SupportedTypes tells the router which targets this module can handle
	SupportedTypes() []string
	Fetch(ctx context.Context, target string) (ProviderResult, error)
}

type ProviderResult struct {
	ProviderName string                 `json:"provider_name"`
	Target       string                 `json:"target"`
	RawData      map[string]interface{} `json:"raw_data"`
	RiskScore    int                    `json:"risk_score"`
}
