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

# resource "jetstream_consumer" "ORDERS_NEW" {
#   stream_id      = jetstream_stream.ORDERS.id
#   durable_name   = "NEW"
#   deliver_all    = true
#   filter_subject = "ORDERS.received"
#   sample_freq    = 100
# }

