output "helm_release" {
  value = module.upi_platform.release_name
}

output "namespace" {
  value = module.upi_platform.namespace
}

output "ingress_host" {
  value = "Add to hosts file: 127.0.0.1 upi.local (or use kubectl port-forward)"
}
