data "k8s_resource" "example" {
  api_version = "v1"
  kind        = "ConfigMap"
  metadata = {
    name      = "test"
    namespace = "default"
  }
}
