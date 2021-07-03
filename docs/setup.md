# Setup

## 1. Install caddy with provider

Download from caddy from [here](https://caddyserver.com/download).

## 2. Put Initialize config

```bash
$ cat /etc/caddy/config.json
{
  "admin": {
    "listen": "0.0.0.0:2019",
    "config": {
      "persist": true
    }
  }
}
```

## 3. Start caddy

```bash
$ sudo caddy run -config /etc/caddy/config.json -resume /var/tmp
```

## 4. Start ingress-external-caddy

```bash
$ ./ingress-external-caddy --kubeconfig ~/.kube/config --provider cloudflare --cloudflare-email "user@example.com" --cloudflare-api-token "dummy-token" --backend "localhost:8080" --caddy-host "caddy-host"
```

## 5. Block outside configuration

caddy will listen `0.0.0.0`. It can change the configuration from the world. So We **force recommend** block HTTP access from the untrust world.

For example...

- iptables
- firewalld
