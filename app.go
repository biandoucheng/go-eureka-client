package goeurekaclient

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// AppObject 应用信息
type AppObject struct {
	Name  string          //应用名称
	Hosts []AddressObject //主机地址
}

// NewApp 实例化应用对象
func NewApp(name string) AppObject {
	return AppObject{
		Name:  name,
		Hosts: []AddressObject{},
	}
}

// AddHost 添加主机实例
func (a *AppObject) AddHost(scheme, host, port, healthUrl string) bool {
	addr := NewAddress(a.Name, scheme, host, port, healthUrl)

	// 检测该地址是否已存在
	for _, a_adr := range a.Hosts {
		if a_adr.Equl(addr) {
			return true
		}
	}

	a.Hosts = append(a.Hosts, addr)
	return true
}

// HasInstance 判断该应用是否存在实例
func (a *AppObject) HasInstance() bool {
	return len(a.Hosts) > 0
}

// HasActiveInstance 判断该应用是否有确定健康的实例
func (a *AppObject) HasActiveInstance() bool {
	a.RefresHost()
	return len(a.Hosts) > 0
}

// RefresHost 刷新应用的实例信息,去除不健康的主机信息
func (a *AppObject) RefresHost() {
	addrs := []AddressObject{}

	for _, addr := range a.Hosts {
		if !addr.Check() {
			continue
		}

		addrs = append(addrs, addr)
	}

	a.Hosts = addrs
}

// GetAnUrl 获取一个地址
func (a *AppObject) GetAnUrl() (string, error) {
	if len(a.Hosts) == 0 {
		return "", errors.New("Get app(" + a.Name + ") address failed with err: no adders cached")
	}

	index := rand.Intn(len(a.Hosts))
	addr := a.Hosts[index]

	return addr.GetUrl(), nil
}
