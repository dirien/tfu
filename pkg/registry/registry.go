package registry

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type RegistryDetails struct {
	ID        string      `json:"id"`
	Owner     string      `json:"owner"`
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Alias     interface{} `json:"alias"`
	Version   string      `json:"version"`
	Versions  []string    `json:"versions"`
}

type RegistryType string

const (
	Providers RegistryType = "providers"
	Modules   RegistryType = "modules"
)

func GetRegistryDetails(provider string, registryType RegistryType) (*RegistryDetails, error) {
	client := resty.New()

	resp, err := client.R().
		Get(fmt.Sprintf("https://registry.terraform.io/v1/%s/%s", registryType, provider))
	if err != nil {
		return nil, err
	}
	registryProvider := RegistryDetails{}
	if err := json.Unmarshal(resp.Body(), &registryProvider); err != nil {
		return nil, err
	}
	return &registryProvider, nil
}
