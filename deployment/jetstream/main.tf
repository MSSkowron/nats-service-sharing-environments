terraform {
  required_providers {
    jetstream = {
      source = "nats-io/jetstream"
      version = "0.0.35"
    }
  }
}

provider "jetstream" {
  servers     = var.servers
  user        = var.login
  password    = var.password
}

resource "jetstream_stream" "ORDERS" {
  name     = var.stream_name
  subjects = var.stream_subjects
  storage  = "file"
  max_age  = 60 * 60 * 24 * 365
}

