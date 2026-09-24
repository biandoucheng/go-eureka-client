package goeurekaclient

import (
	"errors"
	"log"
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
	log.Printf("apps: %+v \n", e.Apps)
	e.L.RUnlock()
}

// GetAnHost
func (e *EurekaAppCache) GetAnHost(cfname string, name string) (AddressObject, error) {
	name = strings.ToUpper(name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name

	e.L.RLock()
	app, ok := e.Apps[kname]
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
	app, ok := e.Apps[kname]
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

// GetAllUrls 获取应用的全部服务地址。
func (e *EurekaAppCache) GetAllUrls(cfname string, name string) ([]string, error) {
	name = strings.ToUpper(name)
	cfname = strings.ToUpper(cfname)
	kname := cfname + "_" + name

	e.L.RLock()
	app, ok := e.Apps[kname]
	e.L.RUnlock()

	if !ok {
		return nil, errors.New("Get app url failed with err: app (" + name + ") not found")
	}

	if !app.HasHost() {
		return nil, errors.New("Get app url failed with err: app (" + name + ") has no address")
	}

	return app.GetAllUrls(), nil
}

// GetAppUrl
func GetAppUrl(cfname string, name string) (string, error) {
	return globalEurekaAppCache.GetAnUrl(cfname, name)
}

// GetAllUrls 获取应用的全部服务地址。
func GetAllUrls(cfname string, name string) ([]string, error) {
	return globalEurekaAppCache.GetAllUrls(cfname, name)
}

// GetAppHost
func GetAnHost(cfname string, name string) (AddressObject, error) {
	return globalEurekaAppCache.GetAnHost(cfname, name)
}

// ShowApps
func ShowApps() {
	globalEurekaAppCache.ShowApps()
}
