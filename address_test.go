package goeurekaclient

import "testing"

func TestAddressEqualityIncludesScheme(t *testing.T) {
	httpAddress := NewAddress("APP", "http", "127.0.0.1", "8080", "/health")
	httpsAddress := NewAddress("APP", "https", "127.0.0.1", "8080", "/health")

	if httpAddress.Equl(httpsAddress) {
		t.Fatal("HTTP and HTTPS addresses should not be considered equal")
	}
}
