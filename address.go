package goeurekaclient

import (
	"strings"
)

// AddressObject 地址对象
type AddressObject struct {
	AppName   string // 应用名称
	Scheme    string // 协议头
	Host      string // 主机地址
	Port      string // 服务端口
	HealthUrl string // 健康接口地址
}

// NewAddress 实例化一个地址对象
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

// Check 检车主机是否可用
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

// Equl 检查两个地址是否一致
func (a *AddressObject) Equl(addr AddressObject) bool {
	return a.Host == addr.Host && a.Port == addr.Port
}

// GetUrl 获取完整的请求地址
func (a *AddressObject) GetUrl() string {
	return a.Scheme + "://" + a.Host + ":" + a.Port
}
