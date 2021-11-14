terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "=3.75.0"
    }
  }
}

provider "aws" {
  version = "~> 3.0"
  region  = "us-west-2"
}