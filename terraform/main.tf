terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.25.0"
    }
  }
}

variable "prefix" {
  type = string
  default = "httptops"
}

variable "message_retention_duration" {
  type = string
  default = "86600s"
}

resource "google_service_account" "main" {
  account_id = "${var.prefix}-publisher"
  description = "Publishes events to PubSub"
}

resource "google_service_account_key" "credentials" {
  service_account_id = google_service_account.main.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

resource "google_pubsub_topic" "main" {
  name = "${var.prefix}-messages"
  message_retention_duration = var.message_retention_duration
}

resource "google_pubsub_topic_iam_member" "member" {
  topic = google_pubsub_topic.main.name
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${google_service_account.main.email}"
}

output "topic" {
  value = google_pubsub_topic.main.name
}

output "service_account_key" {
  value = google_service_account_key.credentials.private_key
}
