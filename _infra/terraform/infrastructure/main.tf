module "api" {
  source          = "../modules/cloudrun"
  stage           = var.stage
  name            = "api-${var.region}"
  project         = var.project
  region          = var.region
  service_account = module.api_service_account.email
  cpu             = var.api_cpu
  memory          = var.api_memory
  min_instances   = var.api_min_instances
  max_instances   = var.api_max_instances
  ingress         = "all"
  image           = "eu.gcr.io/${var.project}/api:latest"
  concurrency     = 100
}
