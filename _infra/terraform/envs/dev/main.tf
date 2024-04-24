provider "google" {
  project = var.project
  region  = var.region
}

terraform {
  required_version = "~>1.5.6"
    backend "gcs" {
      bucket = "waracle-test-dev-terraform"
    }
}

module "infrastructure" {
  source  = "../../infrastructure"
  stage   = var.stage
  project = var.project
  region = var.region
}