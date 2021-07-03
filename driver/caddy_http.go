package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"k8s.io/klog/v2"
)

// UpdateConfig request http POST to /load
func UpdateConfig(conf caddy.Config, host string) error {
	u := fmt.Sprintf("http://%s:2019/load", host)

	d, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(d))
	if err != nil {
		return fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	req.Header.Add("Cache-Control", "must-revalidate")
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to HTTP POST: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("invalid status code (code: %d, body: %s)", resp.StatusCode, b)
	}

	klog.Info("Update successfully!")

	return nil
}
