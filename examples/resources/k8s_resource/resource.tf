resource "k8s_resource" "example" {
  manifest = {
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata = {
      name      = "example"
      namespace = "default"
    }
    data = {
      foo = "bar"
    }
  }
}
