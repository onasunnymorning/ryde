package ryde

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

// TestIsValidType tests the IsValidType function.
// It tests if the function returns the expected boolean value for different input types.
func TestIsValidType(t *testing.T) {
	tests := []struct {
		name string
		t    string
		want bool
	}{
		{
			name: "valid type FULL",
			t:    "FULL",
			want: true,
		},
		{
			name: "valid type DIFF",
			t:    "DIFF",
			want: true,
		},
		{
			name: "invalid type",
			t:    "INVALID",
			want: false,
		},
		{
			name: "mixed case type",
			t:    "fUlL",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidType(tt.t); got != tt.want {
				t.Errorf("IsValidType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewXMLDeposit tests the NewXMLDeposit function which creates a new XMLDeposit object.
// It tests the function with different input parameters and checks if the output is as expected.
func TestNewXMLDeposit(t *testing.T) {
	watermark := time.Now()
	tests := []struct {
		name    string
		t       string
		id      string
		prevId  string
		resend  int
		wantErr error
	}{
		{
			name:    "valid FULL type",
			t:       "FULL",
			id:      "123",
			prevId:  "456",
			resend:  0,
			wantErr: nil,
		},
		{
			name:    "valid DIFF type",
			t:       "DIFF",
			id:      "789",
			prevId:  "012",
			resend:  1,
			wantErr: nil,
		},
		{
			name:    "invalid type",
			t:       "INVALID",
			id:      "345",
			prevId:  "678",
			resend:  0,
			wantErr: ErrInvalidDepositType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep, err := NewXMLDeposit(tt.t, tt.id, tt.prevId, tt.resend, watermark)
			if err != tt.wantErr {
				t.Errorf("NewXMLDeposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Only test the fields if there is no error, otherwise the test will panic.
			if err == nil {
				if dep.Type != strings.ToUpper(tt.t) {
					t.Errorf("NewXMLDeposit() Type = %v, want %v", dep.Type, strings.ToUpper(tt.t))
				}
				if dep.ID != tt.id {
					t.Errorf("NewXMLDeposit() ID = %v, want %v", dep.ID, tt.id)
				}
				if dep.PrevID != tt.prevId {
					t.Errorf("NewXMLDeposit() PrevID = %v, want %v", dep.PrevID, tt.prevId)
				}
				if dep.Resend != tt.resend {
					t.Errorf("NewXMLDeposit() Resend = %v, want %v", dep.Resend, tt.resend)
				}
				if dep.Watermark != watermark.Format(time.RFC3339) {
					t.Errorf("NewXMLDeposit() Watermark = %v, want %v", dep.Watermark, watermark.Format(time.RFC3339))
				}
				if dep.Domain != NameSpace["domain"] {
					t.Errorf("NewXMLDeposit() Domain = %v, want %v", dep.Domain, NameSpace["domain"])
				}
				if dep.Contact != NameSpace["contact"] {
					t.Errorf("NewXMLDeposit() Contact = %v, want %v", dep.Contact, NameSpace["contact"])
				}
				if dep.SecDNS != NameSpace["secDNS"] {
					t.Errorf("NewXMLDeposit() SecDNS = %v, want %v", dep.SecDNS, NameSpace["secDNS"])
				}
				if dep.Rde != NameSpace["rde"] {
					t.Errorf("NewXMLDeposit() Rde = %v, want %v", dep.Rde, NameSpace["rde"])
				}
				if dep.RdeHeader != NameSpace["rdeHeader"] {
					t.Errorf("NewXMLDeposit() RdeHeader = %v, want %v", dep.RdeHeader, NameSpace["rdeHeader"])
				}
				if dep.RdeDomain != NameSpace["rdeDomain"] {
					t.Errorf("NewXMLDeposit() RdeDomain = %v, want %v", dep.RdeDomain, NameSpace["rdeDomain"])
				}
				if dep.RdeHost != NameSpace["rdeHost"] {
					t.Errorf("NewXMLDeposit() RdeHost = %v, want %v", dep.RdeHost, NameSpace["rdeHost"])
				}
				if dep.RdeContact != NameSpace["rdeContact"] {
					t.Errorf("NewXMLDeposit() RdeContact = %v, want %v", dep.RdeContact, NameSpace["rdeContact"])
				}
				if dep.RdeRegistrar != NameSpace["rdeRegistrar"] {
					t.Errorf("NewXMLDeposit() RdeRegistrar = %v, want %v", dep.RdeRegistrar, NameSpace["rdeRegistrar"])
				}
				if dep.RdeIDN != NameSpace["rdeIDN"] {
					t.Errorf("NewXMLDeposit() RdeIDN = %v, want %v", dep.RdeIDN, NameSpace["rdeIDN"])
				}
				if dep.RdeNNDN != NameSpace["rdeNNDN"] {
					t.Errorf("NewXMLDeposit() RdeNNDN = %v, want %v", dep.RdeNNDN, NameSpace["rdeNNDN"])
				}
				if dep.RdeEppParams != NameSpace["rdeEppParams"] {
					t.Errorf("NewXMLDeposit() RdeEppParams = %v, want %v", dep.RdeEppParams, NameSpace["rdeEppParams"])
				}
				if dep.RdePolicy != NameSpace["rdePolicy"] {
					t.Errorf("NewXMLDeposit() RdePolicy = %v, want %v", dep.RdePolicy, NameSpace["rdePolicy"])
				}
				if dep.Epp != NameSpace["epp"] {
					t.Errorf("NewXMLDeposit() Epp = %v, want %v", dep.Epp, NameSpace["epp"])
				}
			}
		})
	}
}

// Test if the XMLDeposit object can Unmarshal a valid XML string
func TestXMLDeposit_UnmarshalXML(t *testing.T) {
	deposit := &XMLDepositUnMarshall{}
	err := xml.Unmarshal([]byte(getValidFullDepositXMLString()), deposit)
	if err != nil {
		t.Errorf("UnmarshalXML() error = %v", err)
	}
	if deposit.Type != "FULL" {
		t.Errorf("UnmarshalXML() Type = %v, want %v", deposit.Type, "FULL")
	}
	if deposit.ID != "20191017001" {
		t.Errorf("UnmarshalXML() ID = %v, want %v", deposit.ID, "20191017001")
	}
	if deposit.PrevID != "" {
		t.Errorf("UnmarshalXML() PrevID = %v, want %v", deposit.PrevID, "")
	}
	if deposit.Resend != 0 {
		t.Errorf("UnmarshalXML() Resend = %v, want %v", deposit.Resend, 0)
	}
	if deposit.Watermark != "2019-10-17T00:00:00Z" {
		t.Errorf("UnmarshalXML() Watermark = %v, want %v", deposit.Watermark, "2019-10-17T00:00:00Z")
	}
}

// Helper function to get a valid XML string for testing purposes
func getValidFullDepositXMLString() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
	<rde:deposit type="FULL" id="20191017001"
	  xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"
	  xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"
	  xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1"
	  xmlns:rde="urn:ietf:params:xml:ns:rde-1.0"
	  xmlns:rdeHeader="urn:ietf:params:xml:ns:rdeHeader-1.0"
	  xmlns:rdeDomain="urn:ietf:params:xml:ns:rdeDomain-1.0"
	  xmlns:rdeHost="urn:ietf:params:xml:ns:rdeHost-1.0"
	  xmlns:rdeContact="urn:ietf:params:xml:ns:rdeContact-1.0"
	  xmlns:rdeRegistrar="urn:ietf:params:xml:ns:rdeRegistrar-1.0"
	  xmlns:rdeIDN="urn:ietf:params:xml:ns:rdeIDN-1.0"
	  xmlns:rdeNNDN="urn:ietf:params:xml:ns:rdeNNDN-1.0"
	  xmlns:rdeEppParams="urn:ietf:params:xml:ns:rdeEppParams-1.0"
	  xmlns:rdePolicy="urn:ietf:params:xml:ns:rdePolicy-1.0"
	  xmlns:epp="urn:ietf:params:xml:ns:epp-1.0">
	
	  <rde:watermark>2019-10-17T00:00:00Z</rde:watermark>
	  <rde:rdeMenu>
		<rde:version>1.0</rde:version>
		<rde:objURI>urn:ietf:params:xml:ns:rdeHeader-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeContact-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeHost-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeDomain-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeRegistrar-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeIDN-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeNNDN-1.0
		</rde:objURI>
		<rde:objURI>urn:ietf:params:xml:ns:rdeEppParams-1.0
		</rde:objURI>
	  </rde:rdeMenu>
	
	  <!-- Contents -->
	  <rde:contents>
		<!-- Header -->
		<rdeHeader:header>
		  <rdeHeader:tld>test</rdeHeader:tld>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeDomain-1.0">2
			</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeHost-1.0">1
			</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeContact-1.0">1
			</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeRegistrar-1.0">1
		</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeIDN-1.0">1
			</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeNNDN-1.0">1
			</rdeHeader:count>
		  <rdeHeader:count
			uri="urn:ietf:params:xml:ns:rdeEppParams-1.0">1
		</rdeHeader:count>
		</rdeHeader:header>
	
		<!-- Domain: example1.example -->
		<rdeDomain:domain>
		  <rdeDomain:name>example1.example</rdeDomain:name>
		  <rdeDomain:roid>Dexample1-TEST</rdeDomain:roid>
		  <rdeDomain:status s="ok"/>
		  <rdeDomain:registrant>jd1234</rdeDomain:registrant>
		  <rdeDomain:contact type="admin">sh8013</rdeDomain:contact>
		  <rdeDomain:contact type="tech">sh8013</rdeDomain:contact>
		  <rdeDomain:ns>
			<domain:hostObj>ns1.example.com</domain:hostObj>
			<domain:hostObj>ns1.example1.example</domain:hostObj>
		  </rdeDomain:ns>
		  <rdeDomain:clID>RegistrarX</rdeDomain:clID>
		  <rdeDomain:crRr client="jdoe">RegistrarX</rdeDomain:crRr>
		  <rdeDomain:crDate>1999-04-03T22:00:00.0Z</rdeDomain:crDate>
		  <rdeDomain:exDate>2025-04-03T22:00:00.0Z</rdeDomain:exDate>
		</rdeDomain:domain>
	
		<!-- Domain: example2.example -->
		<rdeDomain:domain>
		  <rdeDomain:name>example2.example</rdeDomain:name>
		  <rdeDomain:roid>Dexample2-TEST</rdeDomain:roid>
		  <rdeDomain:status s="ok"/>
		  <rdeDomain:status s="clientUpdateProhibited"/>
		  <rdeDomain:registrant>jd1234</rdeDomain:registrant>
		  <rdeDomain:contact type="admin">sh8013</rdeDomain:contact>
		  <rdeDomain:contact type="tech">sh8013</rdeDomain:contact>
		  <rdeDomain:clID>RegistrarX</rdeDomain:clID>
		  <rdeDomain:crRr>RegistrarX</rdeDomain:crRr>
		  <rdeDomain:crDate>1999-04-03T22:00:00.0Z</rdeDomain:crDate>
		  <rdeDomain:exDate>2025-04-03T22:00:00.0Z</rdeDomain:exDate>
		</rdeDomain:domain>
	
		<!-- Host: ns1.example.example -->
		<rdeHost:host>
		  <rdeHost:name>ns1.example1.example</rdeHost:name>
		  <rdeHost:roid>Hns1_example_test-TEST</rdeHost:roid>
		  <rdeHost:status s="ok"/>
		  <rdeHost:status s="linked"/>
		  <rdeHost:addr ip="v4">192.0.2.2</rdeHost:addr>
		  <rdeHost:addr ip="v4">192.0.2.29</rdeHost:addr>
		  <rdeHost:addr ip="v6">2001:DB8:1::1</rdeHost:addr>
		  <rdeHost:clID>RegistrarX</rdeHost:clID>
		  <rdeHost:crRr>RegistrarX</rdeHost:crRr>
		  <rdeHost:crDate>1999-05-08T12:10:00.0Z</rdeHost:crDate>
		  <rdeHost:upRr>RegistrarX</rdeHost:upRr>
		  <rdeHost:upDate>2009-10-03T09:34:00.0Z</rdeHost:upDate>
		</rdeHost:host>
	
		<!-- Contact: sh8013 -->
		<rdeContact:contact>
		  <rdeContact:id>sh8013</rdeContact:id>
		  <rdeContact:roid>Csh8013-TEST</rdeContact:roid>
		  <rdeContact:status s="linked"/>
		  <rdeContact:status s="clientDeleteProhibited"/>
		  <rdeContact:postalInfo type="int">
			<contact:name>John Doe</contact:name>
			<contact:org>Example Inc.</contact:org>
			<contact:addr>
			  <contact:street>123 Example Dr.</contact:street>
			  <contact:street>Suite 100</contact:street>
			  <contact:city>Dulles</contact:city>
			  <contact:sp>VA</contact:sp>
			  <contact:pc>20166-6503</contact:pc>
			  <contact:cc>US</contact:cc>
			</contact:addr>
		  </rdeContact:postalInfo>
		  <rdeContact:voice x="1234">+1.7035555555
		  </rdeContact:voice>
		  <rdeContact:fax>+1.7035555556
		  </rdeContact:fax>
		  <rdeContact:email>jdoe@example.example
		  </rdeContact:email>
		  <rdeContact:clID>RegistrarX</rdeContact:clID>
		  <rdeContact:crRr client="jdoe">RegistrarX
		  </rdeContact:crRr>
		  <rdeContact:crDate>2009-09-13T08:01:00.0Z
		  </rdeContact:crDate>
		  <rdeContact:upRr client="jdoe">RegistrarX
		  </rdeContact:upRr>
		  <rdeContact:upDate>2009-11-26T09:10:00.0Z
		  </rdeContact:upDate>
		  <rdeContact:trDate>2009-12-03T09:05:00.0Z
		  </rdeContact:trDate>
		  <rdeContact:disclose flag="0">
			<contact:voice/>
			<contact:email/>
		  </rdeContact:disclose>
		</rdeContact:contact>
	
		<!-- Registrar: RegistrarX -->
		<rdeRegistrar:registrar>
		  <rdeRegistrar:id>RegistrarX</rdeRegistrar:id>
		  <rdeRegistrar:name>Registrar X</rdeRegistrar:name>
		  <rdeRegistrar:gurid>8</rdeRegistrar:gurid>
		  <rdeRegistrar:status>ok</rdeRegistrar:status>
		  <rdeRegistrar:postalInfo type="int">
			<rdeRegistrar:addr>
			  <rdeRegistrar:street>123 Example Dr.
			  </rdeRegistrar:street>
			  <rdeRegistrar:street>Suite 100
			  </rdeRegistrar:street>
			  <rdeRegistrar:city>Dulles</rdeRegistrar:city>
			  <rdeRegistrar:sp>VA</rdeRegistrar:sp>
			  <rdeRegistrar:pc>20166-6503</rdeRegistrar:pc>
			  <rdeRegistrar:cc>US</rdeRegistrar:cc>
			</rdeRegistrar:addr>
		  </rdeRegistrar:postalInfo>
		  <rdeRegistrar:voice x="1234">+1.7035555555
		  </rdeRegistrar:voice>
		  <rdeRegistrar:fax>+1.7035555556
		  </rdeRegistrar:fax>
		  <rdeRegistrar:email>jdoe@example.example
		  </rdeRegistrar:email>
		  <rdeRegistrar:url>http://www.example.example
		  </rdeRegistrar:url>
		  <rdeRegistrar:whoisInfo>
			<rdeRegistrar:name>whois.example.example
			</rdeRegistrar:name>
			<rdeRegistrar:url>http://whois.example.example
			</rdeRegistrar:url>
		  </rdeRegistrar:whoisInfo>
		  <rdeRegistrar:crDate>2005-04-23T11:49:00.0Z
		  </rdeRegistrar:crDate>
		  <rdeRegistrar:upDate>2009-02-17T17:51:00.0Z
		  </rdeRegistrar:upDate>
		</rdeRegistrar:registrar>
	
		<!-- IDN Table -->
		<rdeIDN:idnTableRef id="pt-BR">
		  <rdeIDN:url>
	http://www.iana.org/domains/idn-tables/tables/br_pt-br_1.0.html
		  </rdeIDN:url>
		  <rdeIDN:urlPolicy>
			http://registro.br/dominio/regras.html
		  </rdeIDN:urlPolicy>
		</rdeIDN:idnTableRef>
	
		<!-- NNDN: pinguino.example -->
		<rdeNNDN:NNDN>
		  <rdeNNDN:aName>xn--exampl-gva.example</rdeNNDN:aName>
		  <rdeNNDN:idnTableId>pt-BR</rdeNNDN:idnTableId>
		  <rdeNNDN:originalName>example1.example</rdeNNDN:originalName>
		  <rdeNNDN:nameState>withheld</rdeNNDN:nameState>
		  <rdeNNDN:crDate>2005-04-23T11:49:00.0Z</rdeNNDN:crDate>
		</rdeNNDN:NNDN>
	
		<!-- EppParams -->
		<rdeEppParams:eppParams>
		  <rdeEppParams:version>1.0</rdeEppParams:version>
		  <rdeEppParams:lang>en</rdeEppParams:lang>
		  <rdeEppParams:objURI>
			urn:ietf:params:xml:ns:domain-1.0
		  </rdeEppParams:objURI>
		  <rdeEppParams:objURI>
			urn:ietf:params:xml:ns:contact-1.0
		  </rdeEppParams:objURI>
		  <rdeEppParams:objURI>
			urn:ietf:params:xml:ns:host-1.0
		  </rdeEppParams:objURI>
		  <rdeEppParams:svcExtension>
			<epp:extURI>urn:ietf:params:xml:ns:rgp-1.0
			</epp:extURI>
			<epp:extURI>urn:ietf:params:xml:ns:secDNS-1.1
			</epp:extURI>
		  </rdeEppParams:svcExtension>
		  <rdeEppParams:dcp>
		  <epp:access><epp:all/></epp:access>
			<epp:statement>
			  <epp:purpose>
				<epp:admin/>
				<epp:prov/>
			  </epp:purpose>
			  <epp:recipient>
				<epp:ours/>
				<epp:public/>
			  </epp:recipient>
			  <epp:retention>
				<epp:stated/>
			  </epp:retention>
			</epp:statement>
		  </rdeEppParams:dcp>
		</rdeEppParams:eppParams>
	  <rdePolicy:policy
		 scope="//rde:deposit/rde:contents/rdeDomain:domain"
		 element="rdeDomain:registrant" />
	  </rde:contents>
	</rde:deposit>`
}
