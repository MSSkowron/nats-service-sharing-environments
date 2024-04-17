variable "name" {
    description = "name of the Helm release"
    type = string
    default = "nats"
}

variable "values" {
    description = "path to a yaml file containing value overrides"
    type = string
    default = "nats/cluster.yaml"
}

variable "kubeconfig" {
    description = "path to kubeconfig"
    type = string
    default = "~/.kube/config"
}

variable "namespace" {
    description = "kubernetes namespace of the NATS cluster deployment"
    type = string
    default = "default"
}