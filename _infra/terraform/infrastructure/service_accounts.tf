module "api_service_account" {
  source        = "terraform-google-modules/service-accounts/google"
  version       = "~> 4.0"
  project_id    = var.project
  names         = ["api-sa"]
  project_roles = [
    "${var.project}=>roles/editor",
    "${var.project}=>roles/secretmanager.secretAccessor",
    "${var.project}=>roles/iam.serviceAccountTokenCreator"
  ]
}