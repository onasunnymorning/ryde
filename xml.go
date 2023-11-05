package ryde

import "encoding/xml"

// https://www.rfc-editor.org/rfc/rfc9022.html#name-xml-model

// Represents an XML escrow deposit tag. Can be used to Marshall and Unmarshall XML deposit element.
type XMLDeposit struct {
	XMLName      xml.Name `xml:"urn:ietf:params:xml:ns:rde-1.0 deposit" json:"-"`
	Type         string   `xml:"type,attr"`
	ID           string   `xml:"id,attr"`
	PrevID       string   `xml:"prevId,attr"`
	Resend       int      `xml:"resend,attr"`
	Watermark    string   `xml:"urn:ietf:params:xml:ns:rde-1.0 watermark"`
	Domain       string   `xml:"xmlns:domain,attr"`
	Contact      string   `xml:"xmlns:contact,attr"`
	SecDNS       string   `xml:"xmlns:secDNS,attr"`
	Rde          string   `xml:"xmlns:rde,attr"`
	RdeHeader    string   `xml:"xmlns:rdeHeader,attr"`
	RdeDomain    string   `xml:"xmlns:rdeDomain,attr"`
	RdeHost      string   `xml:"xmlns:rdeHost,attr"`
	RdeContact   string   `xml:"xmlns:rdeContact,attr"`
	RdeRegistrar string   `xml:"xmlns:rdeRegistrar,attr"`
	RdeIDN       string   `xml:"xmlns:rdeIDN,attr"`
	RdeNNDN      string   `xml:"xmlns:rdeNNDN,attr"`
	RdeEppParams string   `xml:"xmlns:rdeEppParams,attr"`
	RdePolicy    string   `xml:"xmlns:rdePolicy,attr"`
	Epp          string   `xml:"xmlns:epp,attr"`
}

// Represents an XML header tag. Can be used to Marshall and Unmarshall XML header element.
type RDEHeader struct {
	TLD       string     `xml:"tld"`
	Registrar int        `xml:"registrar"`
	PPSP      int        `xml:"ppsp"`
	Count     []RDECount `xml:"count"`
}

type RDECount struct {
	Uri string `xml:"uri,attr" json:"object"`
	ID  int    `xml:",chardata" json:"count"`
}
