package goeurekaclient

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// AppObject
type AppObject struct {
	Name  string
	Hosts []AddressObject
}

// NewApp
func NewApp(name string) AppObject {
	return AppObject{
		Name:  name,
		Hosts: []AddressObject{},
	}
}

// AddHost
func (a *AppObject) AddHost(scheme, host, port, healthUrl string) bool {
	addr := NewAddress(a.Name, scheme, host, port, healthUrl)

	for _, a_adr := range a.Hosts {
		if a_adr.Equl(addr) {
			return true
		}
	}

	a.Hosts = append(a.Hosts, addr)
	return true
}

// HasHost
func (a *AppObject) HasHost() bool {
	return len(a.Hosts) > 0
}

// GetAddresses
func (a *AppObject) GetAddresses() []AddressObject {
	return a.Hosts
}

// RemoveUnhealthAddress
func (a *AppObject) RemoveUnhealthAddress(addrs []AddressObject) {
	if len(addrs) == 0 {
		return
	}

	unhealths := map[string]bool{}
	for _, adr := range addrs {
		unhealths[adr.Url()] = false
	}

	healths := []AddressObject{}
	for _, adr := range a.Hosts {
		if _, has := unhealths[adr.Url()]; !has {
			healths = append(healths, adr)
		}
	}

	a.Hosts = healths
}

// GetAnUrl
func (a *AppObject) GetAnUrl() (string, error) {
	if len(a.Hosts) == 0 {
		return "", errors.New("Get app(" + a.Name + ") address failed with err: no adders cached")
	}

	index := rand.Intn(len(a.Hosts))
	addr := a.Hosts[index]

	return addr.Url(), nil
}

// GetAllUrls
func (a *AppObject) GetAllUrls() []string {
	if len(a.Hosts) == 0 {
		return []string{}
	}

	uls := []string{}
	for _, adr := range a.Hosts {
		uls = append(uls, adr.Url())
	}

	return uls
}

// GetAnHost
func (a *AppObject) GetAnHost() (AddressObject, error) {
	if len(a.Hosts) == 0 {
		return AddressObject{}, errors.New("Get app(" + a.Name + ") address failed with err: no adders cached")
	}

	index := rand.Intn(len(a.Hosts))
	addr := a.Hosts[index]

	return addr, nil
}
