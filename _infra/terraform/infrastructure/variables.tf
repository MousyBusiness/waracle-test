variable "stage" {}

variable "project" {}

variable "region" {}

variable "api_cpu" {
  default = "1"    # 1 vCPUs
}
variable "api_memory" {
  default = "256Mi"  # 256MB memory
}
variable "api_min_instances" {
  default = 0
}
variable "api_max_instances" {
  default = 1
}