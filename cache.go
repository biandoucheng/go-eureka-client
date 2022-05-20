package goeurekaclient

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

// EurekaAppCache Eureka服务下应用列表缓存
type EurekaAppCache struct {
	L    sync.RWMutex         // 读写锁
	Apps map[string]AppObject // 维护的应用信息列表
}

// 定义全局应用维护列表
var globalEurekaAppCache = EurekaAppCache{
	Apps: map[string]AppObject{},
}

// Save 存储应用信息
func (e *EurekaAppCache) Save(info EurekaAppInfo) {
	name := strings.ToUpper(info.Name)
	app, ok := e.Apps[name]
	if !ok {
		app = NewApp(name)
	}

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

	defer e.L.Unlock()
	e.L.Lock()
	e.Apps[name] = app
}

// ClearAdderss 清除应用下无用的地址
// 慎用该方法
// 该方法是为了保证本地缓存应用的健康情况
// 所以使用时确保在服务端注册的应用的健康接口可以Get调通
func (e *EurekaAppCache) ClearAdderss(name string) {
	app, ok := e.Apps[name]
	if !ok {
		return
	}

	app.RefresHost()
	defer e.L.Unlock()
	e.L.Lock()
	e.Apps[name] = app
}

// ClearUseless 清除无效应用实例(健康接口不通)
// 慎用该方法
// 该方法是为了保证本地缓存应用的健康情况
// 所以使用时确保在服务端注册的应用的健康接口可以Get调通
func (e *EurekaAppCache) ClearUseless() {
	app_ks := make([]string, len(e.Apps))
	index := 0

	for key := range e.Apps {
		app_ks[index] = key
		index += 1
	}

	for _, name := range app_ks {
		e.ClearAdderss(name)
	}
}

// GetAppUrl 获取一个应用的请求地址
func GetAppUrl(name string) (string, error) {
	name = strings.ToUpper(name)

	app, ok := globalEurekaAppCache.Apps[name]
	if !ok {
		return "", errors.New("Get app url failed with err: app (" + name + ") not found")
	}

	ul, err := app.GetAnUrl()
	if err != nil {
		return "", err
	}

	return ul, nil
}
