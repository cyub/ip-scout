package units

import (
	"errors"
	"fmt"
	"net"

	geoip2 "github.com/oschwald/geoip2-golang"
)

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

func (geoip *GeoIP) City(ip string) error {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return err
	}

	defer db.Close()

	userIP := net.ParseIP(ip)
	if userIP == nil {
		err := fmt.Sprintf("Invalid Ip:%v", userIP)
		return errors.New(err)
	}

	record, err := db.City(userIP)
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
