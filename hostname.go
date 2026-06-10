package main

import (
	"net"
	"strings"
)

// ResolveHost performs a reverse DNS lookup on an IP address.
// If a hostname is found, it returns the hostname without the
// trailing period. If no hostname exists, it returns the IP address.
func ResolveHost(ip string) string {

	//Reverse DNS lookup
	names, err := net.LookupAddr(ip)

	//Return the first hostname if one is available
	if err == nil && len(names) > 0 {

		//Remove the trailing "." commonly returned by DNS
		return strings.TrimSuffix(names[0], ".")
	}

	// Fallback to the original IP address when no hostname exists
	return ip
}
