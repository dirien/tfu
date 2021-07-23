package registry

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Provider struct {
	ID        string      `json:"id"`
	Owner     string      `json:"owner"`
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Alias     interface{} `json:"alias"`
	Version   string      `json:"version"`
	Versions  []string    `json:"versions"`
}

func GetRegistryProvider(provider string) (*Provider, error) {
	client := resty.New()

	resp, err := client.R().
		Get("https://registry.terraform.io/v1/providers/" + provider)
	if err != nil {
		return nil, err
	}
	registryProvider := Provider{}
	if err := json.Unmarshal(resp.Body(), &registryProvider); err != nil {
		return nil, err
	}
	return &registryProvider, nil
}
