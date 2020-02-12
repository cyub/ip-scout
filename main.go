package main

import (
	"encoding/json"
	"fmt"
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
		if userIP == "" {
			userIP = RemoteAddr(r).(string)
		}

		if env != "production" {
			log.Printf("the client ip is %v", userIP)
		}

		userAgent := r.Header.Get("User-Agent")
		var isCurlRequest bool = false
		if userAgent[:4] == "curl" {
			isCurlRequest = true
		}
		err := geoip.City(userIP)
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			result := "Opps! 系统异常啦"
			w.Header().Set("Content-Length", strconv.Itoa(len(result)))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(result))
			fmt.Println(err)
			return
		}

		if isCurlRequest != true {
			data, _ := json.Marshal(geoip)
			w.Header().Set("Content-Type", "application/json")
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

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(result)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
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
	ips := proxyIps(r)
	if len(ips) > 0 && ips[0] != "" {
		rip, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			rip = ips[0]
		}
		return rip
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil
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
