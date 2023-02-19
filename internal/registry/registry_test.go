package registry

import (
	"fmt"
	"os"

	"github.com/dirien/tfu/internal/hcl"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestProviderEquals(t *testing.T) {
	d := RegistryDetailsProvider{
		client: resty.New(),
	}

	provider := hcl.Terraform{
		RequiredProviders: hcl.RequiredProviders{
			Providers: map[string]map[string]string{
				"google": map[string]string{
					"source":  "hashicorp/google",
					"version": "3.75.0",
				},
			},
		},
	}

	result := &RegistryDetails{
		ID:        "hashicorp/google/4.0.0",
		Owner:     "hashicorp",
		Namespace: "hashicorp",
		Name:      "google",
		Alias:     "google",
		Version:   "4.0.0",
		Versions: []string{
			"0.1.0",
			"4.0.0",
		},
	}

	httpmock.ActivateNonDefault(d.client.GetClient())
	defer httpmock.DeactivateAndReset()
	r, err := os.ReadFile("testdata/google.provider.json")
	if err != nil {
		require.NoError(t, err)
	}
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://registry.terraform.io/v1/%s/%s", Providers, provider.RequiredProviders.Providers["google"]["source"]),
		httpmock.NewStringResponder(200, string(r)))

	details, err := d.GetRegistryDetails(provider.RequiredProviders.Providers["google"]["source"], Providers)
	if err != nil {
		require.NoError(t, err)
	}
	require.Equal(t, details, result)
}
