package providers

import (
	"context"
	"fmt"
	"net"
	"strings"
)

type DomainDNSSecurityProvider struct {
	Resolver *net.Resolver
}

func NewDomainDNSSecurityProvider() *DomainDNSSecurityProvider {
	return &DomainDNSSecurityProvider{
		Resolver: &net.Resolver{},
	}
}

func (p *DomainDNSSecurityProvider) Name() string {
	return "domain_dns_security"
}

func (p *DomainDNSSecurityProvider) SupportedTypes() []string {
	return []string{"DOMAIN"}
}

func (p *DomainDNSSecurityProvider) Fetch(ctx context.Context, target string) (ProviderResult, error) {
	risk := 0
	spfRecord := ""
	dmarcRecord := ""
	status := "secure"

	// 1. Check SPF on the root domain
	txtRecords, err := p.Resolver.LookupTXT(ctx, target)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				spfRecord = txt
				break
			}
		}
	}

	// 2. Check DMARC on the _dmarc subdomain
	dmarcTarget := fmt.Sprintf("_dmarc.%s", target)
	dmarcTxt, err := p.Resolver.LookupTXT(ctx, dmarcTarget)
	if err == nil {
		for _, txt := range dmarcTxt {
			if strings.HasPrefix(txt, "v=DMARC1") {
				dmarcRecord = txt
				break
			}
		}
	}

	// --- RISK CALCULATION ALGORITHM ---
	if spfRecord == "" {
		risk += 20 // Missing SPF allows basic spoofing
	}

	if dmarcRecord == "" {
		risk += 30 // Missing DMARC is a critical underwriting failure
		status = "vulnerable_to_spoofing"
	} else {
		// DMARC exists, but is it enforced?
		if strings.Contains(dmarcRecord, "p=none") {
			risk += 15 // p=none is monitoring only, not preventing spoofing
			status = "dmarc_monitoring_only"
		} else if strings.Contains(dmarcRecord, "p=quarantine") || strings.Contains(dmarcRecord, "p=reject") {
			// Properly enforced
			status = "enforced"
		}
	}

	if risk > 100 {
		risk = 100
	}

	rawData := map[string]interface{}{
		"spf_record":   spfRecord,
		"dmarc_record": dmarcRecord,
		"status":       status,
	}

	return ProviderResult{
		ProviderName: p.Name(),
		Target:       target,
		RawData:      rawData,
		RiskScore:    risk,
	}, nil
}
