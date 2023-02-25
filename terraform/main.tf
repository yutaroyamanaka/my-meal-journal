terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.55.0"
    }
  }
  required_version = ">= 1.3.9"

  cloud {
    organization = "my-meal-journal"

    workspaces {
      name = "ci"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
