variable "kubeconfig_path" {
  type    = string
  default = "~/.kube/config"
}

variable "kube_context" {
  type    = string
  default = "kind-upi-local"
}

variable "chart_path" {
  type    = string
  default = "../../../helm/upi-platform"
}
