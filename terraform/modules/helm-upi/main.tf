resource "helm_release" "upi_platform" {
  name             = var.release_name
  chart            = var.chart_path
  namespace        = var.namespace
  create_namespace = true
  wait             = var.wait
  timeout          = var.timeout

  values = [for f in var.values_files : file(f)]

  dynamic "set" {
    for_each = var.set_values
    content {
      name  = set.key
      value = set.value
    }
  }
}
