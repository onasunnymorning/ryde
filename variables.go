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
		"eppParams":           EPP_PARAMS_FILE_SUFFIX,
		"policy":              POLICY_FILE_SUFFIX,
		// TODO: move this to a separate part in the struct as its not a csv file
		"analysis": ANALYSYS_FILE_SUFFIX,
	}

	CSVFilesAndHeaders = map[string][]string{
		"domain":              RDE_DOMAIN_CSV_HEADER,
		"domainStatus":        RDE_DOMAIN_STATUS_CSV_HEADER,
		"domainNameservers":   RDE_DOMAIN_NAMESERVER_CSV_HEADER,
		"domainDnssec":        RDE_DOMAIN_DNSSEC_CSV_HEADER,
		"domainTransfers":     RDE_DOMAIN_TRANSFER_CSV_HEADER,
		"host":                RDE_HOST_CSV_HEADER,
		"hostAddress":         RDE_HOST_ADDRESS_CSV_HEADER,
		"hostStatus":          RDE_HOST_STATUS_CSV_HEADER,
		"contact":             RDE_CONTACT_CSV_HEADER,
		"contactStatus":       RDE_CONTACT_STATUS_CSV_HEADER,
		"contactPostalInfo":   RDE_CONTACT_PINFO_CSV_HEADER,
		"registrar":           RDE_REGISTRAR_CSV_HEADER,
		"registrarPostalInfo": RDE_REGISTRAR_PINFO_CSV_HEADER,
		"idnLanguage":         RDE_IDNLANGUAGE_CSV_HEADER,
		"nndn":                RDE_NNDN_CSV_HEADER,
		"uniqueContactID":     RDE_UNIQUE_CONTACT_ID_CSV_HEADER,
		"eppParams":           RDE_EPP_PARAMS_CSV_HEADER,
		"policy":              RDE_POLICY_CSV_HEADER,
	}

	// These define the order of the columns in the CSV files and the number of columns is used to validate the CSV files
	RDE_DOMAIN_CSV_HEADER             = []string{"Name", "RoID", "Uname", "IdnTableId", "OriginalName", "Registrant", "ClID", "CrRr", "CrDate", "ExDate", "UpRr", "UpDate"}
	RDE_DOMAIN_STATUS_CSV_HEADER      = []string{"Name", "Status"}
	RDE_DOMAIN_NAMESERVER_CSV_HEADER  = []string{"Name", "HostObjID"}
	RDE_DOMAIN_DNSSEC_CSV_HEADER      = []string{"Name", "KeyTag", "Algorithm", "DigestType", "Digest"}
	RDE_DOMAIN_TRANSFER_CSV_HEADER    = []string{"Name", "TrStatus", "ReID", "ReDate", "AcID", "AcDate", "ExDate"}
	RDE_HOST_CSV_HEADER               = []string{"Name", "RoID", "ClID", "CrRr", "CrDate", "UpRr", "UpDate"}
	RDE_HOST_ADDRESS_CSV_HEADER       = []string{"Name", "Version", "Address"}
	RDE_HOST_STATUS_CSV_HEADER        = []string{"Name", "Status"}
	RDE_NNDN_CSV_HEADER               = []string{"AName", "UName", "IDNTableID", "OriginalName", "NameState", "CrDate"}
	RDE_REGISTRAR_MAPPING_CSV_HEARDER = []string{"ClID", "Name", "GurID", "RegistrarID", "DomainCount", "HostCount", "ContactCount"}
	RDE_CONTACT_CSV_HEADER            = []string{"ID", "RoID", "Voice", "Fax", "Email", "ClID", "CrRr", "CrDate", "UpRr", "UpDate"}
	RDE_CONTACT_STATUS_CSV_HEADER     = []string{"ID", "Status"}
	RDE_CONTACT_PINFO_CSV_HEADER      = []string{"ID", "Type", "Name", "Org", "Street1", "Street2", "Street3", "City", "SP", "PC", "CC"}
	RDE_REGISTRAR_CSV_HEADER          = []string{"ID", "Name", "GurID", "status", "WhoisURL", "URL", "CrDate", "UpDate", "Voice", "Fax", "email"}
	RDE_REGISTRAR_PINFO_CSV_HEADER    = []string{"ID", "Type", "Street1", "Street2", "Street3", "City", "StateProvince", "PostalCode", "CountryCode"}
	RDE_IDNLANGUAGE_CSV_HEADER        = []string{"ID", "URL", "URLPolicy"}
	RDE_UNIQUE_CONTACT_ID_CSV_HEADER  = []string{"ID"}
	RDE_POLICY_CSV_HEADER             = []string{"Scope", "Element"}
	RDE_EPP_PARAMS_CSV_HEADER         = []string{"Version", "Lang", "Value"}

	CounterValidationRules = []CounterValidationRule{
		{"domain", "rdeDomain", true},
		{"host", "rdeHost", true},
		{"contact", "rdeContact", true},
		{"registrar", "rdeRegistrar", true},
		{"idnLanguage", "rdeIDN", true},
		{"nndn", "rdeNNDN", true},
		{"eppParams", "rdeEppParams", true},
		{"domainStatus", "rdeDomain", false},
		{"domainNameservers", "rdeDomain", false},
		{"domainDnssec", "rdeDomain", false},
		{"domainTransfers", "rdeDomain", false},
		{"hostAddress", "rdeHost", false},
		{"hostStatus", "rdeHost", false},
		{"contactStatus", "rdeContact", false},
		{"contactPostalInfo", "rdeContact", false},
		{"registrarPostalInfo", "rdeRegistrar", false},
		{"uniqueContactID", "rdeContact", false},
		{"policy", "rdePolicy", true},
	}
)
