package ryde

import "encoding/xml"

type XMLNNDN struct {
	XMLName      xml.Name `xml:"NNDN"`
	AName        string   `xml:"aName"`
	UName        string   `xml:"uName"`
	IDNTableID   string   `xml:"idnTableId"`
	OriginalName string   `xml:"originalName"`
	NameState    string   `xml:"nameState"`
	CrDate       string   `xml:"crDate"`
}
