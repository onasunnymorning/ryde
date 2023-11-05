package ryde

import (
	"encoding/xml"
	"strings"
	"time"
)

// Represents an XML escrow deposit tag. Can be used to Marshall XML deposit element.
// https://www.rfc-editor.org/rfc/rfc9022.html#name-xml-model
type XMLDepositMarshall struct {
	XMLName      xml.Name `xml:"urn:ietf:params:xml:ns:rde-1.0 rde:deposit" json:"-"`
	Type         string   `xml:"type,attr"`
	ID           string   `xml:"id,attr"`
	PrevID       string   `xml:"prevId,attr"`
	Resend       int      `xml:"resend,attr"`
	Watermark    string   `xml:"rde:watermark"`
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

// NewXMLDeposit creates a new XMLDeposit object with the given parameters.
// If the deposit type is invalid, it returns an error.
// The returned object contains various namespaces for different elements.
func NewXMLDeposit(t, id, prevId string, resend int, watermark time.Time) (*XMLDepositMarshall, error) {
	if !IsValidType(t) {
		return nil, ErrInvalidDepositType
	}
	return &XMLDepositMarshall{
		Type:         strings.ToUpper(t),
		ID:           id,
		PrevID:       prevId,
		Resend:       resend,
		Watermark:    watermark.Format(time.RFC3339),
		Domain:       NameSpace["domain"],
		Contact:      NameSpace["contact"],
		SecDNS:       NameSpace["secDNS"],
		Rde:          NameSpace["rde"],
		RdeHeader:    NameSpace["rdeHeader"],
		RdeDomain:    NameSpace["rdeDomain"],
		RdeHost:      NameSpace["rdeHost"],
		RdeContact:   NameSpace["rdeContact"],
		RdeRegistrar: NameSpace["rdeRegistrar"],
		RdeIDN:       NameSpace["rdeIDN"],
		RdeNNDN:      NameSpace["rdeNNDN"],
		RdeEppParams: NameSpace["rdeEppParams"],
		RdePolicy:    NameSpace["rdePolicy"],
		Epp:          NameSpace["epp"],
	}, nil
}

// Represents an XML escrow deposit tag. Can be used to UnMarshall XML deposit element.
// https://stackoverflow.com/questions/48609596/xml-namespace-prefix-issue-at-go
type XMLDepositUnMarshall struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:rde-1.0 deposit" json:"-"`
	Type      string   `xml:"type,attr"`
	ID        string   `xml:"id,attr"`
	PrevID    string   `xml:"prevId,attr"`
	Resend    int      `xml:"resend,attr"`
	Watermark string   `xml:"watermark"`
}

// IsValidType checks if the given string is a valid type.
// A valid type is either "FULL" or "DIFF".
func IsValidType(t string) bool {
	return strings.ToUpper(t) == "FULL" || strings.ToUpper(t) == "DIFF"
}
