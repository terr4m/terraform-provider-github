data "k8s_resources" "example" {
  api_version = "v1"
  kind        = "ConfigMap"
}
