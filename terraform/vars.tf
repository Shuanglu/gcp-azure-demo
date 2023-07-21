variable "gcp_sa_id" {
  type = string
}

variable "gcp_issuer" {
  type    = string
  default = "https://accounts.google.com"
}

variable "azure_sub_id" {
  type = string
}
