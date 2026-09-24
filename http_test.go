package goeurekaclient

import (
	"io"
	"net/http"
	"testing"
)

type trackingBody struct {
	closed bool
}

func (b *trackingBody) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSetHTTPTransport(t *testing.T) {
	original := getHTTPTransport()
	t.Cleanup(func() {
		SetHTTPTransport(original)
	})

	called := false
	SetHTTPTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(http.NoBody),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	}))

	resp, err := HttpGet("http://eureka.test/apps", nil, nil, 1)
	if err != nil {
		t.Fatalf("HttpGet() error = %v", err)
	}
	defer resp.Body.Close()

	if !called {
		t.Fatal("configured HTTP transport was not used")
	}
}

func TestSetHTTPTransportNilRestoresDefault(t *testing.T) {
	SetHTTPTransport(roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, nil
	}))

	SetHTTPTransport(nil)
	if _, ok := getHTTPTransport().(*http.Transport); !ok {
		t.Fatal("nil transport did not restore default *http.Transport")
	}
}

func TestEurekaResponseBodyClosed(t *testing.T) {
	original := getHTTPTransport()
	t.Cleanup(func() {
		SetHTTPTransport(original)
	})

	body := &trackingBody{}
	SetHTTPTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       body,
			Header:     make(http.Header),
			Request:    req,
		}, nil
	}))

	if _, err := EurekaGetApp("http://eureka.test", "", "APP"); err == nil {
		t.Fatal("EurekaGetApp() error = nil, want HTTP status error")
	}
	if !body.closed {
		t.Fatal("EurekaGetApp() did not close response body")
	}
}
