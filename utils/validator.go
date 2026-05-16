package utils

import (
	"net"
	"regexp"
	"strings"
)

type TargetType string

const (
	TargetTypeEmail   TargetType = "EMAIL"
	TargetTypeDomain  TargetType = "DOMAIN"
	TargetTypeIP      TargetType = "IP"
	TargetTypeUnknown TargetType = "UNKNOWN"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IdentifyTarget analyzes the string and classifies it
func IdentifyTarget(target string) TargetType {
	target = strings.TrimSpace(target)

	// Check if it's an IP address
	if net.ParseIP(target) != nil {
		return TargetTypeIP
	}

	// Check if it's an email
	if emailRegex.MatchString(target) {
		return TargetTypeEmail
	}

	// Simple heuristic for domains (contains a dot, no spaces)
	if strings.Contains(target, ".") && !strings.Contains(target, " ") {
		return TargetTypeDomain
	}

	return TargetTypeUnknown
}
