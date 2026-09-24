package goeurekaclient

import (
	"errors"
	"net"
	"time"
)

// GetInnerIp 获取内网IP
func GetInnerIp() (string, error) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("get local interface addresses failed: " + err.Error())
	}

	for _, item := range ips {
		if ipnet, ok := item.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("get local IPv4 address failed")
}

// GetMs 获取毫秒时间戳
func GetMs() int64 {
	return time.Now().UnixNano() / 1e6
}
