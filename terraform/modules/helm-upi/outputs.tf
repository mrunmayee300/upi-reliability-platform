output "release_name" {
  value = helm_release.upi_platform.name
}

output "namespace" {
  value = helm_release.upi_platform.namespace
}

output "status" {
  value = helm_release.upi_platform.status
}
