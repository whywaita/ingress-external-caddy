# ingress-external-caddy

ingress-external-caddy configure an external [Caddy](https://caddyserver.com/) from [Kubernetes ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/).

```text
 +-------------------+                     +-------------------+ 
 | Kubernetes        |                     | caddy in internet |
 |                   |   -- Configure ->   |                   |
 | ingress           |                     | host              | <- Internet
 |  - a.example.com  | <- reverse proxy -- |  - a.example.com  |
 |  - b.example.com  |                     |  - b.example.com  |
 +-------------------+                     +-------------------+ 
```

## Options

| options | description |
|:---:|:---:|
| `--kubeconfig` | file path of kubeconfig |
| `--backend` | URL of backend. caddy will set to upstream |
| `--provider` | provider of ACME |
| `--caddy-host` | IP or hostname of caddy provisioned |
| `--cloudflare-email` | Email address for cloudflare |
| `--cloudflare-token` | API Token for cloudflare |

### Support provider

- `cloudflare`: [DNS-01 cloudflare](https://certbot-dns-cloudflare.readthedocs.io/en/latest/)