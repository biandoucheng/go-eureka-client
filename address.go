package goeurekaclient

import (
	"strings"
)

// AddressObject
type AddressObject struct {
	AppName   string
	Scheme    string
	Host      string
	Port      string
	HealthUrl string
}

// NewAddress
func NewAddress(name, scheme, host, port, healthUrl string) AddressObject {
	addr := AddressObject{
		AppName:   name,
		Scheme:    scheme,
		Host:      host,
		Port:      port,
		HealthUrl: healthUrl,
	}

	if len(addr.Port) == 0 {
		addr.Port = "80"
	}

	if len(addr.Scheme) == 0 {
		addr.Scheme = "http"
	}

	if !strings.Contains(addr.HealthUrl, addr.Host) && strings.HasPrefix(addr.HealthUrl, "/") {
		addr.HealthUrl = addr.Scheme + "://" + addr.Host + ":" + addr.Port + "/" + strings.TrimLeft(addr.HealthUrl, "/")
	}

	return addr
}

// Check
func (a *AddressObject) Check() bool {
	resp, err := HttpGet(a.HealthUrl, nil, nil, 2)
	if err != nil {
		return false
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return false
	}

	return true
}

// Equl
func (a *AddressObject) Equl(addr AddressObject) bool {
	return a.Host == addr.Host && a.Port == addr.Port
}

// Url
func (a *AddressObject) Url() string {
	return a.Scheme + "://" + a.Host + ":" + a.Port
}
