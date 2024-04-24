// !!! PUBLIC !!!
resource "google_storage_bucket" "public_bucket" {
  name          = "${var.project}-public"
  location      = "eu"
  force_destroy = false
}

data "google_iam_policy" "viewer" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
}

resource "google_storage_bucket_iam_policy" "policy" {
  bucket = google_storage_bucket.public_bucket.name
  policy_data = data.google_iam_policy.viewer.policy_data
}

