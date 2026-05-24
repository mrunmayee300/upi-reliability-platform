variable "release_name" {
  type    = string
  default = "upi-platform"
}

variable "namespace" {
  type    = string
  default = "upi-platform"
}

variable "chart_path" {
  type        = string
  description = "Path to helm/upi-platform chart"
}

variable "values_files" {
  type    = list(string)
  default = []
}

variable "set_values" {
  type    = map(string)
  default = {}
}

variable "wait" {
  type    = bool
  default = true
}

variable "timeout" {
  type    = number
  default = 600
}
