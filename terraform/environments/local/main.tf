provider "kubernetes" {
  config_path    = pathexpand(var.kubeconfig_path)
  config_context = var.kube_context
}

provider "helm" {
  kubernetes {
    config_path    = pathexpand(var.kubeconfig_path)
    config_context = var.kube_context
  }
}

module "upi_platform" {
  source = "../../modules/helm-upi"

  chart_path = abspath(var.chart_path)
  values_files = [
    abspath("${var.chart_path}/values.yaml"),
    abspath("${var.chart_path}/values-local.yaml"),
  ]

  set_values = {
    "global.imagePullPolicy" = "IfNotPresent"
  }
}
