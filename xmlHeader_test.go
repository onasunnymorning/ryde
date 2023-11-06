package ryde

import (
	"encoding/xml"
	"testing"
)

// TestUnMarshallXMLHeader tests the UnMarshallXMLHeader function which unmarshals an XML header string into an XMLHeaderUnMarshall struct.
// It checks if the TLD and count values are correctly unmarshalled and have the expected values.
func TestUnMarshallXMLHeader(t *testing.T) {
	header := XMLHeaderUnMarshall{}
	err := xml.Unmarshal([]byte(getValidHeaderXMLString()), &header)
	if err != nil {
		t.Errorf("Unmarshal failed with error: %v", err)
	}
	if header.TLD != "test" {
		t.Errorf("Expected TLD to be 'test', got '%s'", header.TLD)
	}
	for _, count := range header.Count {
		switch count.Uri {
		case "urn:ietf:params:xml:ns:rdeDomain-1.0":
			if count.ID != 2 {
				t.Errorf("Expected domain count to be 2, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeHost-1.0":
			if count.ID != 1 {
				t.Errorf("Expected host count to be 1, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeContact-1.0":
			if count.ID != 1 {
				t.Errorf("Expected contact count to be 1, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeRegistrar-1.0":
			if count.ID != 1 {
				t.Errorf("Expected registrar count to be 1, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeIDN-1.0":
			if count.ID != 1 {
				t.Errorf("Expected IDN count to be 1, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeNNDN-1.0":
			if count.ID != 1 {
				t.Errorf("Expected NNDN count to be 1, got %d", count.ID)
			}
		case "urn:ietf:params:xml:ns:rdeEppParams-1.0":
			if count.ID != 1 {
				t.Errorf("Expected EPP params count to be 1, got %d", count.ID)
			}
		default:
			t.Errorf("Unexpected URI '%s'", count.Uri)
		}
	}
}

// Returns a valid XML string for the header element.
func getValidHeaderXMLString() string {
	return `<rdeHeader:header>
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
  </rdeHeader:header>`
}
