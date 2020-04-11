package main

import (
	"encoding/json"
	"ip-scout/units"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultAppEnv  = "production"
	defaultAppPort = "8000"
)

var (
	env   = units.Env("APP_ENV", defaultAppEnv)
	port  = units.Env("APP_PORT", defaultAppPort)
	geoip = units.GeoIP{}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userIP := r.URL.Query().Get("ip")
		if len(userIP) == 0 {
			userIP = units.RemoteAddr(r)
		}

		if len(userIP) == 0 {
			writeHTMLToResponse(w, "ip argument invalid", http.StatusBadRequest)
			return
		}

		if env != "production" {
			log.Printf("the client ip is %v", userIP)
		}

		if geoip.City(userIP) != nil {
			writeHTMLToResponse(w, "Opps! 系统异常啦", http.StatusServiceUnavailable)
			return
		}

		isCurlRequest := false
		userAgent := r.Header.Get("User-Agent")
		if len(userAgent) >= 4 && userAgent[:4] == "curl" {
			isCurlRequest = true
		}

		if !isCurlRequest {
			writeJSONToResponse(w, geoip, http.StatusOK)
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
		if geoip.IsAnonymousProxy {
			result += "是\n"
		} else {
			result += "否\n"
		}

		writeHTMLToResponse(w, result, http.StatusOK)
	})
	startWebServer(port)
}

// StartWebServer for start web server
func startWebServer(port string) {
	log.Println("Starting HTTP service at " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("An error occured starting HTTP listener at port %s %s", port, err.Error())
	}
}

func writeJSONToResponse(w http.ResponseWriter, payload interface{}, status int) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func writeHTMLToResponse(w http.ResponseWriter, content string, status int) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(status)
	w.Write([]byte(content))
}
