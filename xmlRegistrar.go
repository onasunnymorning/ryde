// WhoisInfo Value Object
package ryde

import (
	"encoding/xml"
)

type XMLWhoisInfo struct {
	XMLName xml.Name `xml:"whoisInfo" json:"-"`
	Name    string   `xml:"name"`
	URL     string   `xml:"url"`
}

type XMLRegistrarPostalInfo struct {
	XMLName xml.Name `xml:"postalInfo" json:"-"`
	Type    string   `xml:"type,attr"`
	Address XMLAddress
}

type XMLAddress struct {
	XMLName       xml.Name `xml:"addr" json:"-"`
	Street        []string `xml:"street"`
	City          string   `xml:"city"`
	StateProvince string   `xml:"sp"`
	PostalCode    string   `xml:"pc"`
	CountryCode   string   `xml:"cc"`
}

// Registrar Entity
type XMLRegistrar struct {
	XMLName    xml.Name                 `xml:"registrar" json:"-"`
	ID         string                   `xml:"id"`
	Name       string                   `xml:"name"`
	GurID      int                      `xml:"gurid"`
	Status     string                   `xml:"status"`
	PostalInfo []XMLRegistrarPostalInfo `xml:"postalInfo"`
	Voice      string                   `xml:"voice"`
	Fax        string                   `xml:"fax"`
	Email      string                   `xml:"email"`
	URL        string                   `xml:"url"`
	WhoisInfo  XMLWhoisInfo
	CrDate     string `xml:"crDate"`
	UpDate     string `xml:"upDate"`
}
