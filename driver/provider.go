package driver

import (
	"fmt"
	"strings"
)

var (
	// SupportedProvider is list of provider
	SupportedProvider = []string{"cloudflare"}
)

// GetProvider validate provider
func GetProvider(inputProvider string) (string, error) {
	for _, p := range SupportedProvider {
		if strings.EqualFold(inputProvider, p) {
			// is supported provider
			return p, nil
		}
	}

	return "", fmt.Errorf("%s is not supported provider", inputProvider)
}
