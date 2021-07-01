package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"github.com/whywaita/ingress-external-caddy/cmd"
)

// GetCloudFlareTLSConfig get caddytls.TLS using Cloudflare and ACME dns-01
func GetCloudFlareTLSConfig(o cmd.Options, domains []string) (*caddytls.TLS, error) {
	if strings.EqualFold(o.CloudFlareEmail, "") || strings.EqualFold(o.CloudFlareAPIToken, "") {
		return nil, fmt.Errorf("must be set --cloudflare-email and --cloudflare-api-token")
	}

	p := []byte(fmt.Sprintf(`{"api_token": "%s", "name": "cloudflare"}`, o.CloudFlareAPIToken))

	issuer := caddytls.ACMEIssuer{
		Email: o.CloudFlareEmail,
		Challenges: &caddytls.ChallengesConfig{
			DNS: &caddytls.DNSChallengeConfig{
				ProviderRaw: json.RawMessage(p),
			},
		},
	}

	tls := &caddytls.TLS{
		Automation: &caddytls.AutomationConfig{
			Policies: []*caddytls.AutomationPolicy{
				{
					Subjects: domains,
					IssuersRaw: []json.RawMessage{
						caddyconfig.JSONModuleObject(issuer, "module", "acme", nil),
					},
				},
			},
		},
	}

	return tls, nil
}
