package whois

import (
	"errors"
	"net"
)

func IsIP(ip string) bool {
	i := net.ParseIP(ip)
	if i == nil {
		return false
	}
	return true
}

func isPrivateIP(ip string) (bool, error) {
	private := false
	if !IsIP(ip) {
		return false, errors.New("invalid IP address")
	}
	_, prefix08, _ := net.ParseCIDR("10.0.0.0/8")
	_, prefix12, _ := net.ParseCIDR("172.16.0.0/12")
	_, prefix16, _ := net.ParseCIDR("192.168.0.0/16")
	private = prefix08.Contains(net.ParseIP(ip)) || prefix12.Contains(net.ParseIP(ip)) || prefix16.Contains(net.ParseIP(ip))

	return private, nil
}

func ResolvDNS(dns string) ([]string, error) {
	var ips []string
	i, err := net.LookupIP(dns)
	if err != nil {
		return nil, err
	}

	for _, i := range i {
		ips = append(ips, i.String())
	}
	return ips, nil
}
