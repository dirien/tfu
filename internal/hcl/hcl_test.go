package hcl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderEquals(t *testing.T) {
	expected := Terraform{
		RequiredProviders: RequiredProviders{
			Providers: map[string]map[string]string{
				"google": {
					"source":  "hashicorp/google",
					"version": "=3.75.0",
				},
			},
		},
	}

	parser, err := NewHCLFileParser("testdata/main.data")
	if err != nil {
		return
	}
	fmt.Println(parser)
	require.Equal(t, parser.Terraform, expected)
}
