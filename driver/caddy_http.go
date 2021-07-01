package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
)

// UpdateConfig request http POST to /load
func UpdateConfig(conf caddy.Config, host string) error {
	u := fmt.Sprintf("http://%s:2019/load", host)

	d, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal: %w", err)
	}

	resp, err := http.Post(u, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return fmt.Errorf("failed to POST request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("invalid status code (code: %d, body: %s)", resp.StatusCode, b)
	}

	return nil
}
