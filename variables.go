package ryde

var (
	NameSpace = map[string]string{
		"domain":       "urn:ietf:params:xml:ns:domain-1.0",
		"contact":      "urn:ietf:params:xml:ns:contact-1.0",
		"secDNS":       "urn:ietf:params:xml:ns:secDNS-1.1",
		"rde":          "urn:ietf:params:xml:ns:rde-1.0",
		"rdeHeader":    "urn:ietf:params:xml:ns:rdeHeader-1.0",
		"rdeDomain":    "urn:ietf:params:xml:ns:rdeDomain-1.0",
		"rdeHost":      "urn:ietf:params:xml:ns:rdeHost-1.0",
		"rdeContact":   "urn:ietf:params:xml:ns:rdeContact-1.0",
		"rdeRegistrar": "urn:ietf:params:xml:ns:rdeRegistrar-1.0",
		"rdeIDN":       "urn:ietf:params:xml:ns:rdeIDN-1.0",
		"rdeNNDN":      "urn:ietf:params:xml:ns:rdeNNDN-1.0",
		"rdeEppParams": "urn:ietf:params:xml:ns:rdeEppParams-1.0",
		"rdePolicy":    "urn:ietf:params:xml:ns:rdePolicy-1.0",
		"epp":          "urn:ietf:params:xml:ns:epp-1.0",
	}

	CSVFilesAndSuffixes = map[string]string{
		"domain":              DOMAIN_FILE_SUFFIX,
		"domainStatus":        DOMAIN_STATUS_FILE_SUFFIX,
		"domainNameservers":   DOMAIN_NAMESERVER_FILE_SUFFIX,
		"domainDnssec":        DOMAIN_DNSSEC_FILE_SUFFIX,
		"domainTransfers":     DOMAIN_TRANSFER_FILE_SUFFIX,
		"host":                HOST_FILE_SUFFIX,
		"hostAddress":         HOST_ADDRESS_FILE_SUFFIX,
		"hostStatus":          HOST_STATUS_FILE_SUFFIX,
		"contact":             CONTACT_FILE_SUFFIX,
		"contactStatus":       CONTACT_STATUS_FILE_SUFFIX,
		"contactPostalInfo":   CONTACT_PINFO_FILE_SUFFIX,
		"registrar":           REGISTRAR_FILE_SUFFIX,
		"registrarPostalInfo": REGISTRAR_PINFO_FILE_SUFFIX,
		"idnLanguage":         IDN_FILE_SUFFIX,
		"nndn":                NNDN_FILE_SUFFIX,
		"uniqueContactID":     UNIQUE_CONTACT_ID_FILE_SUFFIX,
		// TODO: move this to a separate part in the struct as its not a csv file7pu
		"analysis": ANALYSYS_FILE_SUFFIX,
	}
)
