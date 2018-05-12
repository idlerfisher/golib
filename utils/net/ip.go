package net

import "net"

func GetIpAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	var ips []string
	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	if len(ips) >= 1 {
		return ips[0]
	}
	return ""
}
