package goeurekaclient

import (
	"log"
	"net"
	"time"
)

// GetInnerIp 获取内网IP
func GetInnerIp() string {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range ips {
		if ipnet, ok := item.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetMs 获取毫秒时间戳
func GetMs() int64 {
	return time.Now().UnixNano() / 1e6
}
