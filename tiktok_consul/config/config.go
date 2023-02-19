package config

import (
	"net"
)

//func GetLocalIpAddress() string {
//	// 服务地址硬编码
//	return "192.168.50.145"
//}

//const ConsulIP string = "10.211.55.5"

const ConsulIP = "124.71.11.189"
const ConsulPort = 8500

func GetLocalIpAddress() (ipv4 string, err error) {
	// 获取所有网卡
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	return
}
