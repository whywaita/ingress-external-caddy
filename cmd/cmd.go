package cmd

import "github.com/spf13/cobra"

// Options is options of CLI
type Options struct {
	KubeConfigPath string
	Provider       string
	BackendURL     string
	CaddyHost      string

	CloudFlareEmail    string
	CloudFlareAPIToken string
}

// New create a cmd
func New(o *Options) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().StringVar(&o.KubeConfigPath, "kubeconfig", o.KubeConfigPath, "absolute path of kubeconfig")
	cmd.Flags().StringVar(&o.Provider, "provider", o.Provider, "provider name of ACME")
	cmd.Flags().StringVar(&o.BackendURL, "backend", o.BackendURL, "URL for backend")
	cmd.Flags().StringVar(&o.CaddyHost, "caddy-host", o.CaddyHost, "hostname or ip of caddy")

	// cloudflare provider
	cmd.Flags().StringVar(&o.CloudFlareEmail, "cloudflare-email", o.CloudFlareEmail, "[cloudflare] Email")
	cmd.Flags().StringVar(&o.CloudFlareAPIToken, "cloudflare-api-token", o.CloudFlareAPIToken, "[cloudflare] API Token")
	return cmd
}
