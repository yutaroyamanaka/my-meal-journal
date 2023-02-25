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

  backend "s3" {
    bucket = "bucket-yutaroyamanaka"
    key    = "terraform"
    region = "ap-northeast-1"
    # this is necessary for production env
    # dynamodb_table = "terraform-state-lock"
  }
}

# this is necessary for production env
/*
resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "terraform-state-lock"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
} */

provider "aws" {
  region = "ap-northeast-1"
}
