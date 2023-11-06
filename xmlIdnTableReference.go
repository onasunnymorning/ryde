package ryde

import "encoding/xml"

type XMLIdnTableReference struct {
	XMLName   xml.Name `xml:"idnTableRef" json:"-"`
	ID        string   `xml:"id,attr"`
	Url       string   `xml:"url"`
	UrlPolicy string   `xml:"urlPolicy"`
}
