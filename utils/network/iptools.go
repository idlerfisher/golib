package network

import (
	"net"
	"fmt"
	"math/big"
)

func GetLocalIpAddr() string {
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

func InetNtoA(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) uint32 {
	ret := big.NewInt(0)
	ip1 := net.ParseIP(ip)
	if ip1 != nil {
		if ip1.To4() != nil {
			ret.SetBytes(ip1)
		}
	}
	return uint32(ret.Int64())
}
