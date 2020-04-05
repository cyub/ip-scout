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
		err error
		ips = proxyIps(r)
	)

	if len(ips) > 0 && ips[0] != "" {
		ip, _, err = net.SplitHostPort(ips[0])
	}

	if (ip == "") || (err != nil) {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
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
