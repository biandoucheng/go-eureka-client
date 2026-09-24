package goeurekaclient

import (
	"reflect"
	"testing"
)

func TestGetAllUrlsReturnsAllAddresses(t *testing.T) {
	cache := EurekaAppCache{Apps: map[string]AppObject{}}
	app := NewApp("APP")
	app.AddHost("http", "127.0.0.1", "8080", "/health")
	app.AddHost("https", "127.0.0.1", "8443", "/health")
	cache.Apps["DEFAULT_APP"] = app

	urls, err := cache.GetAllUrls("default", "app")
	if err != nil {
		t.Fatalf("GetAllUrls() error = %v", err)
	}

	want := []string{"http://127.0.0.1:8080", "https://127.0.0.1:8443"}
	if !reflect.DeepEqual(urls, want) {
		t.Fatalf("GetAllUrls() = %v, want %v", urls, want)
	}
}
