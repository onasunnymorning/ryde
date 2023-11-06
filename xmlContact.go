package ryde

import (
	"encoding/xml"
)

type XMLContact struct {
	XMLName    xml.Name               `xml:"contact" json:"-"`
	ID         string                 `xml:"id"`
	RoID       string                 `xml:"roid"`
	Status     []XMLContactStatus     `xml:"status"`
	PostalInfo []XMLContactPostalInfo `xml:"postalInfo"`
	Voice      string                 `xml:"voice"`
	Fax        string                 `xml:"fax"`
	Email      string                 `xml:"email"`
	ClID       string                 `xml:"clID"`
	CrRr       string                 `xml:"crRr"`
	CrDate     string                 `xml:"crDate"`
	UpRr       string                 `xml:"upRr"`
	UpDate     string                 `xml:"upDate"`
	Disclose   XMLDisclose            `xml:"disclose"`
}

type XMLContactPostalInfo struct {
	XMLName xml.Name `xml:"postalInfo" json:"-"`
	Name    string   `xml:"name"`
	Type    string   `xml:"type,attr"`
	Org     string   `xml:"org"`
	Address XMLAddress
}

type XMLContactStatus struct {
	S string `xml:"s,attr"`
}

type XMLDisclose struct {
	Flag  bool                 `xml:"flag,attr"`
	Name  []XMLContactWithType `xml:"name"`
	Org   []XMLContactWithType `xml:"org"`
	Addr  []XMLContactWithType `xml:"addr"`
	Voice []string             `xml:"voice"`
	Fax   []string             `xml:"fax"`
	Email []string             `xml:"email"`
}

type XMLContactWithType struct {
	Type string `xml:"type,attr"`
}
