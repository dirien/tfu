terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "=3.90.1"
    }
  }
}

provider "aws" {
  version = "~> 2.7"
  region  = "us-west-2"
}