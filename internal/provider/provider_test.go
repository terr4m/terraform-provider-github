package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"github": providerserver.NewProtocol6WithError(New("test", "test")()),
}

func testAccPreCheck(t *testing.T) {
}
