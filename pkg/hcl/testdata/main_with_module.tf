terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "= 3.75.0"
    }
  }
}

provider "aws" {
  version = "~> 3.0"
  region  = "us-west-2"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "=3.7.0"
  # insert the 21 required variables here
}

module "instance_replacement_advanced" {
  source = "git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0"

  cloudwatch_log_retention = 14                # Set custom retention for Lambda logs
  name                     = "MY-ASGIR"        # Set custom name
  schedule                 = "rate(5 minutes)" # Set custom check frequency
  timeout                  = "120"             # Set custom timeout
}