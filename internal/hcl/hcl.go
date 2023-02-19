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

provider "aws" {
  version = "~> 2.7"
  region  = "us-west-2"
}

module "instance_replacement_advanced" {
  source = "git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0"

  cloudwatch_log_retention = 14                # Set custom retention for Lambda logs
  name                     = "MY-ASGIR"        # Set custom name
  schedule                 = "rate(5 minutes)" # Set custom check frequency
  timeout                  = "120"             # Set custom timeout
}
*/

// TFFile is the structure of the terraform file
type TFFile struct {
	Terraform Terraform  `hcl:"terraform,block"`
	Provider  []Provider `hcl:"provider,block"`
	Module    []Module   `hcl:"module,block"`
}

type Module struct {
	Name    string `hcl:"name,label"`
	Source  string `hcl:"source"`
	Version string `hcl:"version"`
}

type Provider struct {
	Name    string `hcl:"name,label"`
	Version string `hcl:"version"`
}

type Terraform struct {
	RequiredProviders RequiredProviders `hcl:"required_providers,block"`
}

type RequiredProviders struct {
	Providers map[string]map[string]string `hcl:",remain"`
}

func NewHCLFileParser(filename string) (*TFFile, error) {
	var tfFile TFFile
	parser := hclparse.NewParser()
	hclFile, diag := parser.ParseHCLFile(filename)

	if diag.HasErrors() {
		return nil, errors.New(diag.Error())
	}
	_ = gohcl.DecodeBody(hclFile.Body, nil, &tfFile)
	return &tfFile, nil
}
