package ip_util

import (
	"fmt"
	"net"
)

func GetMyIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		err = fmt.Errorf("net.InterfaceAddrs: %v", err)
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.To4().String()
				break
			}
		}
	}
	return
}
