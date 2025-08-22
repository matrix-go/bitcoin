package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

var kPattern = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHost string, myPort int, startIP int, endIP int, startPort int, endPort int) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	m := kPattern.FindStringSubmatch(address)
	if m == nil {
		return nil
	}

	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])

	neighbors := make([]string, 0)
	for port := startPort; port <= endPort; port++ {
		for ip := startIP; ip < endIP; ip++ {
			guessHost := fmt.Sprintf("%s:%d", prefixHost, lastIp+ip)
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}

func IsFoundHost(host string, port int) bool {
	fmt.Println(host, port)
	return true
}

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil || len(addrs) == 0 {
		return "127.0.0.1"
	}
	return addrs[0]
}
