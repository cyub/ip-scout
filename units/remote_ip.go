package units

import (
	"net"
	"net/http"
	"strings"
)

// RemoteAddr return the client's ip address
func RemoteAddr(r *http.Request) string {
	var (
		ip  string
		ips = proxyIps(r)
	)

	if len(ips) > 0 && ips[0] != "" {
		ip = ips[0]
	} else {
		ip = r.RemoteAddr
		ip, _, _ = net.SplitHostPort(ip)
	}

	if net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

func proxyIps(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}
