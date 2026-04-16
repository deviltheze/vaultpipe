package vault

import (
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// NewClient creates an authenticated Vault client using the given address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address
	cfg.HttpClient = &http.Client{Timeout: 10 * time.Second}

	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	c.SetToken(token)

	return &Client{logical: c.Logical()}, nil
}

// ReadSecrets reads key/value secrets from the given path.
// Supports both KV v1 and KV v2 (detects v2 by "data" wrapper).
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading path %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data := secret.Data

	// KV v2 wraps values under a "data" key.
	if nested, ok := data["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			data = nestedMap
		}
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}

	return result, nil
}
