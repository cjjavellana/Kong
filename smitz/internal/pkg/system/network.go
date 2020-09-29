package system

import (
	"errors"
	"net"
	"runtime"
	"strings"
)

func HostIP() (string, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// handle err
	for _, networkInterface := range networkInterfaces {
		addrs, err := networkInterface.Addrs()
		if err != nil {
			return "", err
		}

		switch runtime.GOOS {
		case "darwin":
			if networkInterface.Name == "en0" {
				return findIPv4(addrs)
			}
		case "linux":
			// TODO: Confirm whether eth1 or eth0
			if networkInterface.Name == "eth0" {
				return findIPv4(addrs)
			}
		}
	}

	return "", errors.New("no ip address found")
}

// Locates the IPv4 ip address
func findIPv4(addrs []net.Addr) (string, error) {
	for _, addr := range addrs {
		// addr.String() returns something like fe80::8dd:363b:b47f:4110/64 for IPv6 and
		// 192.168.1.127/24 for IPv4
		subnet := strings.Split(addr.String(), "/")[1]
		if subnet == "24" {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// process IP address
			if ip != nil {
				return ip.String(), nil
			}
		}

	}

	return "", errors.New("no ipv4 address found")
}
