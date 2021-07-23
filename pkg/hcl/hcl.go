package hcl

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/pkg/errors"
)

/*
	terraform {
	  required_providers {
		google = {
		  source = "hashicorp/google"
		  version = "3.75.0"
		}
	  }
	}
*/
type TFFile struct {
	Terraform Terraform `hcl:"terraform,block"`
}
type Terraform struct {
	RequiredProviders RequiredProviders `hcl:"required_providers,block"`
}

type RequiredProviders struct {
	Providers map[string]map[string]string `hcl:",remain"`
}

func NewHCLRequiredProvidersParser(filename string) (*RequiredProviders, error) {
	var tfFile TFFile
	parser := hclparse.NewParser()
	hclfile, diag := parser.ParseHCLFile(filename)
	if diag.HasErrors() {
		return nil, errors.New(diag.Error())
	}
	_ = gohcl.DecodeBody(hclfile.Body, nil, &tfFile)
	return &tfFile.Terraform.RequiredProviders, nil
}
