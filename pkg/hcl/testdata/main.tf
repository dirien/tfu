terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "=4.1.0"
    }
  }
}

provider "aws" {
  version = "~> 2.7"
  region  = "us-west-2"
}