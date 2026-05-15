package providers

import (
	"context"
	"time"
)

type MockBreachProvider struct{}

func NewMockBreachProvider() *MockBreachProvider {
	return &MockBreachProvider{}
}

func (m *MockBreachProvider) Name() string {
	return "mock_breach_db"
}

func (m *MockBreachProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	// Simulate network delay
	time.Sleep(300 * time.Millisecond)

	data := map[string]interface{}{
		"status": "breached",
		"count":  2,
	}

	return ProviderResult{
		ProviderName: m.Name(),
		Target:       target,
		RawData:      data,
		RiskScore:    85, // Hardcoded for the mock
	}, nil
}