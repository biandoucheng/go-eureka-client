package goeurekaclient

import (
	"strconv"
)

// EurekaClientConfig Eureka客户端配置项
type EurekaClientConfig struct {
	EurekaServerAddress    string   // Eureka服务端接口地址
	Authorization          string   // Http Auth授权信息
	Apps                   []string // 需要的服务名列表
	AppName                string   // 本服务名称
	InstanceDomain         string   // 本服务的域名地址 | 置空
	InstanceIp             string   // 本服务的ip地址 | 置空
	InstancePort           int      // 本服务的开放端口
	InstanceHomePageUrl    string   // 本服务的主页地址
	InstanceStatusUrl      string   // 本服务的状态检查地址
	InstanceHealthCheckUrl string   // 本服务的健康检查地址
	RenewalIntervalInSecs  int64    // 本服务的心跳周期 单位秒
	DurationInSecs         int64    // 本服务的心跳失约后,注册信息保留时长,超时删除注册信息 单位秒
	AppRefreshSecs         int64    // 需要的应用列表里的应用服务信息刷新间隔 单位秒
}

// Id 生成实例ID
func (e *EurekaClientConfig) Id() string {
	return e.InstanceIp + ":" + e.AppName + ":" + strconv.Itoa(e.InstancePort)
}

// HostName 获取主机名称
func (e *EurekaClientConfig) HostName() string {
	if len(e.InstanceDomain) == 0 {
		e.InstanceDomain = e.InstanceIp
	}
	return e.InstanceDomain
}

// RefreshLocalIp 刷新本地IP信息
func (e *EurekaClientConfig) RefreshLocalIp() {
	oip := e.InstanceIp
	e.InstanceIp = GetInnerIp()
	if len(e.InstanceDomain) == 0 || e.InstanceDomain == oip {
		e.InstanceDomain = e.InstanceIp
	}
}

// GlobalEurekaClientConf 全局Eureka客户端配置
var DefaultEurekaClientConf = EurekaClientConfig{
	RenewalIntervalInSecs:  20,
	DurationInSecs:         40,
	AppRefreshSecs:         30,
	InstanceHealthCheckUrl: "/health",
}

// NewEurekaConf 实例化一个eureka客户端配置
func NewEurekaConf() EurekaClientConfig {
	return EurekaClientConfig{
		RenewalIntervalInSecs:  20,
		DurationInSecs:         40,
		AppRefreshSecs:         30,
		InstanceHealthCheckUrl: "/health",
		InstanceIp:             GetInnerIp(),
	}
}

func init() {
	// 默认使用内网IP
	DefaultEurekaClientConf.InstanceIp = GetInnerIp()
}
