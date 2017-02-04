package resolver

import (
	"errors"
	"net"
)

func Address(address string) (ips []string, error error) {
	ips, error = resolveIP(address)

	if ips != nil {
		return
	}

	ips, error = resolveDomain(address)

	if error != nil {
		return
	}

	ips, error = resolveCIDR(address)

	return
}

func resolveIP(address string) (ips []string, error error) {
	ip := net.ParseIP(address)

	if ip != nil {
		ips = []string{ip.String()}
	} else {
		error = errors.New("Invalid address")
	}

	return
}

func resolveDomain(address string) (ips []string, error error) {
	ipCollection, error := net.LookupIP(address)

	if error != nil {
		return
	}

	for _, ip := range ipCollection {
		ips = append(ips, ip.String())
	}

	return
}

func resolveCIDR(address string) (ips []string, error error) {
	ip, network, error := net.ParseCIDR(address)

	if error != nil {
		return
	}

	ip = ip.Mask(network.Mask)

	for network.Contains(ip) {
		for x := len(ip) - 1; x >= 0; x-- {
			ip[x]++

			ips = append(ips, ip.String())

			if ip[x] > 0 {
				break
			}
		}
	}

	return
}
