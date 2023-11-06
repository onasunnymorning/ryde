package ryde

type XMLHost struct {
	Name   string          `xml:"name"`
	RoID   string          `xml:"roid"`
	Status []XMLHostStatus `xml:"status"`
	Addr   []XMLHostAddr   `xml:"addr"`
	ClID   string          `xml:"clID"`
	CrRr   string          `xml:"crRr"`
	CrDate string          `xml:"crDate"`
	UpRr   string          `xml:"upRr"`
	UpDate string          `xml:"upDate"`
}

type XMLHostStatus struct {
	S string `xml:"s,attr"`
}

type XMLHostAddr struct {
	IP string `xml:"ip,attr"`
	ID string `xml:",chardata"`
}
