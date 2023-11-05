package ryde

import "encoding/xml"

// Represents an XML escrow header tag. Can be used to Marshall XML header element.
// https://www.rfc-editor.org/rfc/rfc9022.html#name-header-object
type XMLHeaderMarshall struct {
	XMLName   xml.Name         `xml:"rdeHeader:header" json:"-"`
	TLD       string           `xml:"tld,attr"`
	Registrar int              `xml:"registrar,attr"`
	PPSP      int              `xml:"ppsp,attr"`
	Count     []HeaderURICount `xml:"count"`
}

// Represents an XML escrow deposit tag. Can be used to UnMarshall XML deposit element.
// https://stackoverflow.com/questions/48609596/xml-namespace-prefix-issue-at-go
type XMLHeaderUnMarshall struct {
	XMLName   xml.Name         `xml:"header" json:"-"`
	TLD       string           `xml:"tld"`
	Registrar int              `xml:"registrar"`
	PPSP      int              `xml:"ppsp"`
	Count     []HeaderURICount `xml:"count"`
}

// RDECount represents a count of objects with a given URI.
type HeaderURICount struct {
	Uri string `xml:"uri,attr" json:"object"` // Uri is the URI of the object.
	ID  int    `xml:",chardata" json:"count"` // ID is the count of the object.
}
