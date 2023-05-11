package goeurekaclient

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// globalEurekaAppCache
var globalEurekaAppCache = EurekaAppCache{
	Apps: map[string]AppObject{},
}

// EurekaAppCache
type EurekaAppCache struct {
	L    sync.RWMutex
	Apps map[string]AppObject
}

// Save
func (e *EurekaAppCache) Save(cfname string, info EurekaAppInfo) {
	if len(info.Instance) == 0 {
		return
	}

	name := strings.ToUpper(info.Name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name
	app := NewApp(name)

	for _, ins := range info.Instance {
		schema := "http"
		port := ins.Port.Value
		host := ins.HostName

		if ins.SecurePort.Enable == "true" {
			schema = "https"
			port = ins.SecurePort.Value
		}

		app.AddHost(schema, host, strconv.Itoa(port), ins.HealthCheckUrl)
	}

	e.L.Lock()
	e.Apps[kname] = app
	e.L.Unlock()
}

// ShowApps
func (e *EurekaAppCache) ShowApps() {
	e.L.RLock()
	fmt.Printf("%v", globalEurekaAppCache.Apps)
	e.L.RUnlock()
}

// GetAnHost
func (e *EurekaAppCache) GetAnHost(cfname string, name string) (AddressObject, error) {
	name = strings.ToUpper(name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name

	e.L.RLock()
	app, ok := globalEurekaAppCache.Apps[kname]
	e.L.RUnlock()
	if !ok {
		return AddressObject{}, errors.New("Get app url failed with err: app (" + name + ") not found")
	}

	adr, err := app.GetAnHost()
	if err != nil {
		return AddressObject{}, err
	}

	return adr, nil
}

// GetAnUrl
func (e *EurekaAppCache) GetAnUrl(cfname string, name string) (string, error) {
	name = strings.ToUpper(name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name

	e.L.RLock()
	app, ok := globalEurekaAppCache.Apps[kname]
	e.L.RUnlock()

	if !ok {
		return "", errors.New("Get app url failed with err: app (" + name + ") not found")
	}

	ul, err := app.GetAnUrl()
	if err != nil {
		ul = ""
	}

	return ul, err
}

// GetAllUrl
func (e *EurekaAppCache) GetAllUrl(cfname string, name string) (string, error) {
	name = strings.ToUpper(name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name

	e.L.RLock()
	app, ok := globalEurekaAppCache.Apps[kname]
	e.L.RUnlock()

	if !ok {
		return "", errors.New("Get app url failed with err: app (" + name + ") not found")
	}

	ul, err := app.GetAnUrl()
	if err != nil {
		ul = ""
	}

	return ul, err
}

// GetAppUrl
func GetAppUrl(cfname string, name string) (string, error) {
	return globalEurekaAppCache.GetAnUrl(cfname, name)
}

// GetAppHost
func GetAnHost(cfname string, name string) (AddressObject, error) {
	return globalEurekaAppCache.GetAnHost(cfname, name)
}

// ShowApps
func ShowApps() {
	globalEurekaAppCache.ShowApps()
}
