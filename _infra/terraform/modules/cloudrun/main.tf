resource "google_cloud_run_service" "main" {
  name     = var.name
  location = var.region

  template {
    spec {
      containers {
        image = var.image
        ports {
          container_port = 8080
        }
        resources {
          limits = {
            cpu = var.cpu
            memory = var.memory
          }
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project
        }
        env {
          name  = "REGION"
          value = var.region
        }
      }
      container_concurrency = var.concurrency
      timeout_seconds       = 300
      service_account_name  = var.service_account
    }
    metadata {
      annotations = {
        "client.knative.dev/user-image": var.image
        "autoscaling.knative.dev/minScale" : var.min_instances
        "autoscaling.knative.dev/maxScale" : var.max_instances
        "run.googleapis.com/startup-cpu-boost": "true"
#        "run.googleapis.com/vpc-access-connector" : "svpc-${var.region}",
#        "run.googleapis.com/vpc-access-egress" : "private-ranges-only"
      }
    }
  }

  metadata {
    namespace   = var.project
    annotations = {
      "client.knative.dev/user-image": var.image
      "run.googleapis.com/ingress" = var.ingress
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }


  autogenerate_revision_name = true
}

data "google_iam_policy" "noauth" {
  count = var.require_authentication ? 0 : 1
  binding {
    role    = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  count    = var.require_authentication ? 0 : 1
  location = google_cloud_run_service.main.location
  project  = google_cloud_run_service.main.project
  service  = google_cloud_run_service.main.name

  policy_data = data.google_iam_policy.noauth[count.index].policy_data
}

#resource "google_compute_region_network_endpoint_group" "cloudrun_neg" {
#  name                  = "${var.name}-neg"
#  network_endpoint_type = "SERVERLESS"
#  region                = var.region
#  cloud_run {
#    service = google_cloud_run_service.default.name
#  }
#}