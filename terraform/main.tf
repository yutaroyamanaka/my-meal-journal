terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.55.0"
    }
  }
  required_version = ">= 1.3.9"

  backend "s3" {
    bucket = "bucket-yutaroyamanaka"
    key    = "terraform"
    region = "ap-northeast-1"
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
