package units

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// GeoIP define struct
type GeoIP struct {
	IP               string
	CountryName      string
	CountryCode      string
	CityName         string
	ContinentName    string
	ContinetCode     string
	TimeZone         string
	Latitude         float64
	Longitude        float64
	PostalCode       string
	IsAnonymousProxy bool
}

// City return the ip's city info
func (geoip *GeoIP) City(ip string) error {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return fmt.Errorf("Invalid Ip:%v", netIP)
	}

	geodb := Env("GEOIP_DB", "GeoLite2-City.mmdb")
	db, err := geoip2.Open(geodb)
	if err != nil {
		return err
	}
	defer db.Close()

	record, err := db.City(netIP)
	if err != nil {
		return err
	}

	geoip.IP = ip
	geoip.CountryName = record.Country.Names["zh-CN"]
	geoip.CityName = record.City.Names["en"]
	geoip.CountryCode = record.Country.IsoCode
	geoip.TimeZone = record.Location.TimeZone
	geoip.Latitude = record.Location.Latitude
	geoip.Longitude = record.Location.Longitude
	geoip.ContinentName = record.Continent.Names["zh-CN"]
	geoip.ContinetCode = record.Continent.Code
	geoip.PostalCode = record.Postal.Code
	geoip.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	return nil
}
