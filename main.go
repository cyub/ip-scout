package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	units "github.com/cyub/ip-scout/units"
)

const (
	defaultAppEnv  = "production"
	defaultAppPort = "8000"
)

func main() {
	var (
		env   = Env("APP_ENV", defaultAppEnv)
		port  = Env("APP_PORT", defaultAppPort)
		geoip = units.GeoIP{}
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userIP := r.URL.Query().Get("ip")
		if len(userIP) > 0 && net.ParseIP(userIP) == nil {
			writeHtmlToResponse(w, "ip argument invalid", http.StatusBadRequest)
			return
		}

		if userIP == "" {
			userIP = RemoteAddr(r).(string)
		}

		if env != "production" {
			log.Printf("the client ip is %v", userIP)
		}

		if geoip.City(userIP) != nil {
			writeHtmlToResponse(w, "Opps! 系统异常啦", http.StatusServiceUnavailable)
			return
		}

		var isCurlRequest bool = false
		if r.Header.Get("User-Agent")[:4] == "curl" {
			isCurlRequest = true
		}

		if isCurlRequest != true {
			data, _ := json.Marshal(geoip)
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		result := "当前IP：" + userIP + "\n"
		result += "\n"
		result += "========地理位置========\n"
		result += "国家：" + geoip.CountryName + "\n"
		result += "国家码：" + geoip.CountryCode + "\n"
		result += "城市：" + geoip.CityName + "\n"
		result += "所属洲：" + geoip.ContinentName + "\n"
		result += "洲代码：" + geoip.ContinetCode + "\n"
		result += "时区：" + geoip.TimeZone + "\n"
		result += "经度：" + strconv.FormatFloat(geoip.Longitude, 'f', -1, 64) + "\n"
		result += "纬度：" + strconv.FormatFloat(geoip.Latitude, 'f', -1, 64) + "\n"
		result += "邮政编码：" + geoip.PostalCode + "\n"

		result += "是否匿名代理："
		if geoip.IsAnonymousProxy == true {
			result += "是\n"
		} else {
			result += "否\n"
		}

		writeHtmlToResponse(w, result, http.StatusOK)
	})

	StartWebServer(port)
}

// StartWebServer for start web server
func StartWebServer(port string) {
	log.Println("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}

func Env(item, fallback string) string {
	e := os.Getenv(item)
	if e == "" {
		return fallback
	}
	return e
}

func RemoteAddr(r *http.Request) interface{} {
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

// Proxy returns proxy client ips slice.
func proxyIps(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func writeHtmlToResponse(w http.ResponseWriter, content string, status int) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(status)
	w.Write([]byte(content))
}
