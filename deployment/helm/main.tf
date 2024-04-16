terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
      version = "2.10.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.21.1"
    }
  }
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig
  }
}

provider "kubernetes" {
  config_path    = var.kubeconfig
}

resource "helm_release" "nats" {
  name  = var.name
  chart = "./nats"
  namespace = var.namespace
  create_namespace = true
  values = [file(var.values)]
}

data "kubernetes_service" "nats-lb" {
  depends_on = [helm_release.nats]
  metadata {
    name = "nats-lb"
    namespace = var.namespace
  }
}
