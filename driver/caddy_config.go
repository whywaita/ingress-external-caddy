package driver

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"github.com/whywaita/ingress-external-caddy/cmd"
	networkingv1 "k8s.io/api/networking/v1"
)

// GenerateCaddy generate config of caddyserver
func GenerateCaddy(ingresses []networkingv1.Ingress, o cmd.Options) (*caddy.Config, error) {
	appsConfig, err := getAppsConfig(ingresses, o.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get apps config: %w", err)
	}

	tlsConfig, err := getTLSConfig(ingresses, o)
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}

	caddyHost := net.JoinHostPort(o.CaddyHost, "2019")

	pTrue := true
	c := &caddy.Config{
		Admin: &caddy.AdminConfig{
			Listen:   caddyHost,
			Disabled: false,
			Config: &caddy.ConfigSettings{
				Persist: &pTrue,
			},
			Origins: []string{
				caddyHost,
			},
			EnforceOrigin: true,
		},
		AppsRaw: caddy.ModuleMap{
			"http": caddyconfig.JSON(appsConfig, nil),
			"tls":  caddyconfig.JSON(tlsConfig, nil),
		},
	}

	return c, nil
}

func getAppsConfig(ingresses []networkingv1.Ingress, backend string) (*caddyhttp.App, error) {
	domains := getDomains(ingresses)

	rph := reverseproxy.Handler{
		Upstreams: reverseproxy.UpstreamPool{
			&reverseproxy.Upstream{
				Dial: backend,
			},
		},
	}

	s := &caddyhttp.Server{
		Listen: []string{":80", ":443"},
		Routes: caddyhttp.RouteList{
			caddyhttp.Route{
				MatcherSetsRaw: caddyhttp.RawMatcherSets{
					caddy.ModuleMap{
						"host": caddyconfig.JSON(caddyhttp.MatchHost(domains), nil),
					},
				},
				HandlersRaw: []json.RawMessage{
					caddyconfig.JSONModuleObject(rph, "handler", "reverse_proxy", nil),
				},
			},
		},
	}

	return &caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"ingress": s,
		},
	}, nil
}

func getTLSConfig(ingresses []networkingv1.Ingress, o cmd.Options) (*caddytls.TLS, error) {
	var tlsConfig *caddytls.TLS
	var err error
	switch o.Provider {
	case "cloudflare":
		tlsConfig, err = GetCloudFlareTLSConfig(o, getDomains(ingresses))
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS config using cloudflare: %w", err)
		}
	default:
		return nil, fmt.Errorf("%s is unsupported yet", o.Provider)
	}

	return tlsConfig, nil
}

func getDomains(ingresses []networkingv1.Ingress) []string {
	var domains []string

	for _, i := range ingresses {
		for _, rule := range i.Spec.Rules {
			domains = append(domains, rule.Host)
		}
	}

	return domains
}
