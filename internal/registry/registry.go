package registry

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Details struct {
	ID        string      `json:"id"`
	Owner     string      `json:"owner"`
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Alias     interface{} `json:"alias"`
	Version   string      `json:"version"`
	Versions  []string    `json:"versions"`
}

type Type string

const (
	Providers Type = "providers"
	Modules   Type = "modules"
)

type DetailsProvider struct {
	client *resty.Client
}

func NewRegistryDetails() *DetailsProvider {
	client := resty.New()
	return &DetailsProvider{
		client: client,
	}
}

func (r *DetailsProvider) GetRegistryDetails(provider string, registryType Type) (*Details, error) {
	resp, err := r.client.R().
		Get(fmt.Sprintf("https://registry.terraform.io/v1/%s/%s", registryType, provider))
	if err != nil {
		return nil, err
	}
	registryProvider := Details{}
	if err := json.Unmarshal(resp.Body(), &registryProvider); err != nil {
		return nil, err
	}
	return &registryProvider, nil
}
