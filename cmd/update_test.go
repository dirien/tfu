package tfu

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFile = `terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.75.0"
    }
  }
}

provider "aws" {
  version = "~ 2.7"
  region  = "us-west-2"
}`

func TestUpdateHCT(t *testing.T) {
	file, err := ioutil.TempFile("", "prefix")
	if err != nil {
		require.NoError(t, err)
	}
	_, err = file.Write([]byte(testFile))
	if err != nil {
		require.NoError(t, err)
	}
	defer os.Remove(file.Name())

	err = updateHCLFile(file.Name(), "3.75.0", "4.0.0", false)
	if err != nil {
		require.NoError(t, err)
	}
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		require.NoError(t, err)
	}
	require.Contains(t, string(b), "4.0.0")
}
