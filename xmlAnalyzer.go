package ryde

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// Defines an struct to hold all assets and information about the XML file being analyzed
type XMLAnalyzer struct {
	XMLFile  XMLFile              `json:"xmlFile"`  // The XML file being analyzed.
	CSVFiles map[string]CSVFile   `json:"csvFiles"` // The CSVs file generated during analysis.
	Deposit  XMLDepositUnMarshall `json:"deposit"`  // The struct for containing the UnMarshalled Deposit info
	Header   XMLHeaderUnMarshall  `json:"header"`   // The struct for containing the UnMarshalled Header info
	Counters map[string]int       `json:"counters"` // Holds counters about the number of objects we encountered during analysis. This should match the numbers in the header as well as the number of lines in the CSV files.
	Errors   []string             `json:"errors"`   // Holds errors encountered during analysis.
}

// CSVFile represents a CSV file with its metadata and read/write functionality.
type CSVFile struct {
	FileName       string      `json:"fileName"`  // The name of the CSV file.
	FileSize       int64       `json:"fileSize"`  // The size of the CSV file in bytes.
	LineCount      int         `json:"lineCount"` // The number of lines in the CSV file.
	fileDescriptor *os.File    `json:"-"`         // The file descriptor for the CSV file.
	CsvWriter      *csv.Writer `json:"-"`         // The CSV writer for the CSV file.
}

// XMLFile represents an XML file with its name and size.
type XMLFile struct {
	FileName string       `json:"fileName"` // The name of the XML file.
	FileSize int64        `json:"fileSize"` // The size of the XML file in bytes.
	osFile   *os.File     `json:"-"`
	Decoder  *xml.Decoder `json:"-"`
}

// CounterValidationRule represents a counter validation rule.
type CounterValidationRule struct {
	FileRef     string
	HeaderRef   string
	HeaderMatch bool
}

// NewXMLAnalyzer creates a new instance of XMLAnalyzer and returns a pointer to it.
// It takes a filename string as input and returns an error if the file cannot be opened or its size cannot be determined.
// The function opens the file, checks its size, and saves the filename and size to the XMLFile field of the XMLAnalyzer struct.
func NewXMLAnalyzer(filename string) (*XMLAnalyzer, error) {
	if !strings.HasSuffix(filename, ".xml") {
		return nil, ErrInvalidDepositFileName
	}
	a := XMLAnalyzer{}
	// Try and open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Check and save the file size
	fi, _ := f.Stat() // We can ignore the error here because we know the file exists
	a.XMLFile.FileSize = fi.Size()
	a.XMLFile.FileName = filename
	log.Printf("Created Analyzer for %s (%d MB)\n", a.XMLFile.FileName, a.XMLFile.FileSize/1024/1024)
	a.CSVFiles = make(map[string]CSVFile)
	// Initialize the counters
	c := make(map[string]int)
	for k := range CSVFilesAndSuffixes {
		c[k] = 0
	}
	a.Counters = c
	return &a, nil
}

// Opens the XMLFile and saves the os.File pointer to the XMLFile.osFile field.
func (a *XMLAnalyzer) OpenXMLFile() error {
	reader, err := os.Open(a.XMLFile.FileName)
	if err != nil {
		return err
	}
	a.XMLFile.osFile = reader
	return nil
}

// CloseXMLFile Closes the XMLFile.osFile and removes the pointers from the XMLFile.osFile and XMLFile.Decoder fields.
func (a *XMLAnalyzer) CloseXMLFile() error {
	err := a.XMLFile.osFile.Close()
	if err != nil {
		// TODO: Should we return an error here or just continue?
		return err
	}
	a.XMLFile.osFile = nil
	a.XMLFile.Decoder = nil
	return nil
}

// CreateXMLDecoder Returns an XML decoder for the XMLFile.
func (a *XMLAnalyzer) CreateXMLDecoder() error {
	if a.XMLFile.osFile == nil {
		// TODO: Should we return an error here or open the file?
		return ErrNoXMLReader
	}
	a.XMLFile.Decoder = xml.NewDecoder(a.XMLFile.osFile)
	return nil
}

// AnalyzeDepositTag returns an <rde:deposit> tag by reading the tokens from the decoder and storing the result in the XMLAnalyzer.Deposit field.
func (a *XMLAnalyzer) AnalyzeDepositTag() error {
	progressbar.Default(-1, "Processing <rde:deposit> tag")
	err := a.OpenXMLFile()
	if err != nil {
		return err
	}

	err = a.CreateXMLDecoder()
	if err != nil {
		return err
	}
	defer a.CloseXMLFile()

	err = a.XMLFile.Decoder.Decode(&a.Deposit)
	if err != nil {
		if err == io.EOF {
			return ErrNoDepositTagInFile
		}
		return fmt.Errorf("error decoding deposit: %s", err)
	}

	// found := false

	// for {
	// 	// Stop the loop if we already found and decoded the deposit tag
	// 	if found {
	// 		break
	// 	}
	// 	// Read the next token
	// 	t, tokenErr := a.XMLFile.Decoder.Token()
	// 	if tokenErr != nil {
	// 		if tokenErr == io.EOF {
	// 			return ErrNoDepositTagInFile
	// 		}
	// 		return fmt.Errorf("error decoding token: %s", tokenErr)
	// 	}

	// 	// Only process start elements of type deposit
	// 	switch se := t.(type) {
	// 	case xml.StartElement:
	// 		if se.Name.Local == "deposit" {
	// 			var d = XMLDepositUnMarshall{}
	// 			if err := a.XMLFile.Decoder.DecodeElement(&d, &se); err != nil {
	// 				return fmt.Errorf("error decoding deposit: %s", tokenErr)
	// 			}
	// 			found = true  // Mark as found so we exit as soon as possible
	// 			a.Deposit = d // Save to our XMLAnalyzer struct
	// 		}
	// 	}

	// }
	fmt.Println("Done")
	return nil
}

// AnalyzeTags Analyzes each tag and handle it according to the type of object contained in the tag.
// Decode the tags we are iterested in and stream the data to the appropriate CSV file while collecting counters for sanity checking post analysis.
func (a *XMLAnalyzer) AnalyzeTags() error {
	pbar := progressbar.Default(-1, "Processing object tags")

	err := a.OpenXMLFile()
	if err != nil {
		return err
	}

	err = a.CreateXMLDecoder()
	if err != nil {
		return err
	}
	defer a.CloseXMLFile()

	err = a.CreateCSVFiles()
	if err != nil {
		return err
	}

	err = a.CreateCSVWriters()
	if err != nil {
		return err
	}

	// Create a map to hold unique contact IDs as found on domains. We will use this later to write to a file.
	uniqueContactIDs := make(map[string]bool)

	// Read the entire file, token by token
	for {
		pbar.Add(1)
		// Read the next token
		t, tokenErr := a.XMLFile.Decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				fmt.Println("Done")
				log.Println("Reached end of file")
				break
			}
			return fmt.Errorf("error decoding token: %s", tokenErr)
		}

		// Depending on the token type and Name.Local we handle it accordingly
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "header":
				err := a.processHeaderTag(&se)
				if err != nil {
					return err
				}
			case "registrar":
				err := a.processRegistrarTag(&se)
				if err != nil {
					return err
				}

			case "idnTableRef":
				err := processIDNTableRefTag(&se, a)
				if err != nil {
					return err
				}

			case "contact":
				err := a.processContactTag(&se)
				if err != nil {
					return err
				}

			case "domain":
				err := processDomainTag(&se, a, uniqueContactIDs)
				if err != nil {
					return err
				}

			case "host":
				err := processHostTag(&se, a)
				if err != nil {
					return err
				}

			case "NNDN":
				err := processNNDNTag(&se, a)
				if err != nil {
					return err
				}

			default:
				// Skip all other tags
				continue
			}
		default:
			// Skip all other token types
			continue
		}
	}
	// Now that all tags have been processed
	// Write the unique contact IDs to the file
	log.Println("Writing unique contact IDs to file")
	for k := range uniqueContactIDs {
		err := a.CSVFiles["uniqueContactID"].CsvWriter.Write(StandardizeStringSlice([]string{k}))
		if err != nil {
			return err
		}
	}
	// Finally, flush the writers and close the files
	err = a.FlushCSVWriters()
	if err != nil {
		return err
	}
	err = a.CloseCSVFiles()
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON Writes content of the XMLAnalyzer struct to a JSON file.
// After analysis, the XMLAnalyzer struct contains all the information about the XML file, CSV files and holds several counters and errors.
func (a *XMLAnalyzer) WriteJSON() error {
	//open the file
	file, err := os.OpenFile(a.CSVFiles["analysis"].FileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	// Write the analysis to the file
	fmt.Println("Writing analysis to file")
	analysisBytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(analysisBytes)
	if err != nil {
		return err
	}
	return nil
}

// CheckValidationRules Checks the length of the CSV files against the number of objects in the header and records errors in the analyzer.
func (a *XMLAnalyzer) CheckValidationRules() {
	for _, r := range CounterValidationRules {
		// Always check for an exact match agains the counters we created during processing
		if a.CSVFiles[r.FileRef].LineCount != a.Counters[r.FileRef] {
			a.Errors = append(a.Errors, fmt.Sprintf("CSV file %s has %d lines, expected %d ", r.FileRef, a.CSVFiles[r.FileRef].LineCount, a.Counters[r.FileRef]))
		}
		// Only check against the header if we expect an exact match with the counters in the header
		if r.HeaderMatch {
			for _, v := range a.Header.Count {
				if v.Uri == NameSpace[r.HeaderRef] {
					if v.ID != a.CSVFiles[r.FileRef].LineCount {
						a.Errors = append(a.Errors, fmt.Sprintf("CSV file %s has %d lines, header says %d objects", r.FileRef, a.CSVFiles[r.FileRef].LineCount, v.ID))
					}
				}
			}
		}
	}

}

// Returns the XMLFile.FileName without the file extension.
func (a *XMLAnalyzer) GetBaseXMLFileName() string {
	return strings.Join(strings.Split(a.XMLFile.FileName, ".")[0:len(strings.Split(a.XMLFile.FileName, "."))-1], ".")
}

// Count the number of lines and save the fileSize for the set of CSV files. Use this to check against the number of objects in the header
func (a *XMLAnalyzer) CountLinesInCSVFilesAndSaveSize() error {
	for k, csvFile := range a.CSVFiles {
		file, err := os.OpenFile(csvFile.FileName, os.O_RDONLY, 0444)
		if err != nil {
			return err
		}
		lineCount, err := CountLines(file)
		if err != nil {
			return err
		}
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		csvFile.FileSize = fi.Size()
		csvFile.LineCount = lineCount
		a.CSVFiles[k] = csvFile
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Count the number of columns in the CSV files and match them agains the headers
func (a *XMLAnalyzer) CheckCSVColumnLength() {
	for k, csvFile := range a.CSVFiles {
		// Skip the analysis file
		if k == "analysis" {
			continue
		}
		file, err := os.Open(csvFile.FileName)
		if err != nil {
			a.Errors = append(a.Errors, fmt.Sprintf("error opening file %s: %s", csvFile.FileName, err))
		}
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = len(CSVFilesAndHeaders[k])
		_, err = reader.ReadAll()
		if err != nil {
			a.Errors = append(a.Errors, fmt.Sprintf("CSV column count mismatch in file %s: %s", csvFile.FileName, err))
		}
		file.Close()
	}
}

// Count the number of lines in a file by looking for \n occurrences. Use this to check against the number of objects in the header
func CountLines(r io.Reader) (int, error) {
	var target []byte = []byte("\n")
	var err error

	count := 0

	buffer := make([]byte, 32*1024)

	for {
		read, err := r.Read(buffer)
		if err != nil {
			break
		}

		count += bytes.Count(buffer[:read], target)
	}

	if err == io.EOF {
		return count, nil
	}

	return count, err
}

// Creates CSV files for holding the data from the XML file.
func (a *XMLAnalyzer) CreateCSVFiles() error {
	for k, v := range CSVFilesAndSuffixes {
		csvFile := CSVFile{
			FileName: a.GetBaseXMLFileName() + v,
		}
		a.CSVFiles[k] = csvFile
	}
	return nil
}

// Create csvWriters for each CSV file.
func (a *XMLAnalyzer) CreateCSVWriters() error {
	for k, v := range a.CSVFiles {
		file, err := os.OpenFile(v.FileName, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		v.fileDescriptor = file
		v.CsvWriter = csv.NewWriter(file)
		a.CSVFiles[k] = v
	}
	return nil
}

// Flush csvWriters for each CSV file.
func (a *XMLAnalyzer) FlushCSVWriters() error {
	for k, v := range a.CSVFiles {
		v.CsvWriter.Flush()
		v.CsvWriter = nil
		a.CSVFiles[k] = v
	}
	return nil
}

// Close file descriptors for each CSV file.
func (a *XMLAnalyzer) CloseCSVFiles() error {
	for _, v := range a.CSVFiles {
		v.fileDescriptor.Close()
		v.fileDescriptor = nil
	}
	return nil
}

// Processes a Header object
func (a *XMLAnalyzer) processHeaderTag(se *xml.StartElement) error {
	var header XMLHeaderUnMarshall
	if err := a.XMLFile.Decoder.DecodeElement(&header, se); err != nil {
		return fmt.Errorf("error decoding header: %s", err)
	}
	a.Header = header
	return nil
}

// Processes a Registrar object
func (a *XMLAnalyzer) processRegistrarTag(se *xml.StartElement) error {
	var registrar XMLRegistrar
	// Found a registrar, add it to the counter for sanity checking
	// Skip registrars that are not in the registrar namespace
	if err := a.XMLFile.Decoder.DecodeElement(&registrar, se); err != nil {
		return fmt.Errorf("error decoding registrar: %s", err)
	}
	a.Counters["registrar"]++
	// Prepare the CSV row and standardise the strings and add them to the CSV
	csvRow := []string{registrar.ID, registrar.Name, strconv.Itoa(registrar.GurID), registrar.Status, registrar.WhoisInfo.URL, registrar.URL, registrar.CrDate, registrar.UpDate, registrar.Voice, registrar.Fax, registrar.Email}
	err := a.CSVFiles["registrar"].CsvWriter.Write(StandardizeStringSlice(csvRow))
	if err != nil {
		return err
	}
	// Write the registrar postalinfo to the registrar postalinfo file
	rPostalInfo := make(map[int][]string)
	for i, postalInfo := range registrar.PostalInfo {
		a.Counters["registrarPostalInfo"]++
		rPostalInfo[i] = append(rPostalInfo[i], registrar.ID)
		rPostalInfo[i] = append(rPostalInfo[i], postalInfo.Type)
		// This is clunky but we need to ensure there are always 3 Street elements for CSV length consistency
		// First add the ones that are there
		rPostalInfo[i] = append(rPostalInfo[i], postalInfo.Address.Street...)
		// Then add empty strings for the ones that are missing.
		// A fully slice of strings with 3 street address lines is 4 elements long
		// So we keep adding empty street strings until we reach a lenght of 4
		for len(rPostalInfo[i]) <= 4 {
			rPostalInfo[i] = append(rPostalInfo[i], "")
		}
		rPostalInfo[i] = append(rPostalInfo[i], StandardizeString(postalInfo.Address.City), StandardizeString(postalInfo.Address.StateProvince), StandardizeString(postalInfo.Address.PostalCode), StandardizeString(postalInfo.Address.CountryCode))
	}
	for _, v := range rPostalInfo {
		err := a.CSVFiles["registrarPostalInfo"].CsvWriter.Write(StandardizeStringSlice(v))
		if err != nil {
			return err
		}
	}
	return nil
}

// Processes a Contact object
func (a *XMLAnalyzer) processContactTag(se *xml.StartElement) error {
	// Found a contact, add it to the counter for sanity checking
	// Skip contact tokens that are not in the contact namespace
	if se.Name.Space != NameSpace["rdeContact"] {
		return nil
	}
	a.Counters["contact"]++
	var contact XMLContact
	if err := a.XMLFile.Decoder.DecodeElement(&contact, se); err != nil {
		return fmt.Errorf("error decoding contact: %s", err)
	}
	// Write the contact to the contact file
	contactRow := []string{contact.ID, contact.RoID, contact.Voice, contact.Fax, contact.Email, contact.ClID, contact.CrRr, contact.CrDate, contact.UpRr, contact.UpDate}
	err := a.CSVFiles["contact"].CsvWriter.Write(StandardizeStringSlice(contactRow))
	if err != nil {
		return err
	}
	// Set Status in statusFile
	cStatuses := []string{contact.ID}
	for _, status := range contact.Status {
		cStatuses = append(cStatuses, status.S)
	}
	for i, s := range cStatuses {
		if i == 0 {
			continue
		}
		a.Counters["contactStatus"]++
		err = a.CSVFiles["contactStatus"].CsvWriter.Write(StandardizeStringSlice([]string{contact.ID, s}))
		if err != nil {
			return err
		}
	}
	// Set postalInfo in postalInfoFile
	cPostalInfo := make(map[int][]string)
	for i, postalInfo := range contact.PostalInfo {
		a.Counters["contactPostalInfo"]++
		cPostalInfo[i] = append(cPostalInfo[i], contact.ID)
		cPostalInfo[i] = append(cPostalInfo[i], postalInfo.Type, StandardizeString(postalInfo.Name), StandardizeString(postalInfo.Org))
		// This is clunky but we need to ensure there are always 3 Street elements for CSV length consistency
		// First add the ones that are there
		cPostalInfo[i] = append(cPostalInfo[i], postalInfo.Address.Street...)
		// Then add empty strings for the ones that are missing.
		// A fully slice of strings with 3 street address lines is 6 elements long
		// So we keep adding empty street strings until we reach a lenght of 6
		for len(cPostalInfo[i]) <= 6 {
			cPostalInfo[i] = append(cPostalInfo[i], "")
		}
		cPostalInfo[i] = append(cPostalInfo[i], postalInfo.Address.City, postalInfo.Address.StateProvince, postalInfo.Address.PostalCode, postalInfo.Address.CountryCode)
	}

	for _, v := range cPostalInfo {
		err := a.CSVFiles["contactPostalInfo"].CsvWriter.Write(StandardizeStringSlice(v))
		if err != nil {
			return err
		}
	}
	return nil
}

// Processes a IDNTableRef object
func processIDNTableRefTag(se *xml.StartElement, a *XMLAnalyzer) error {
	// Found an IDN table ref, add it to the counter for sanity checking
	a.Counters["idnLanguage"]++
	var idnTableRef XMLIdnTableReference
	if err := a.XMLFile.Decoder.DecodeElement(&idnTableRef, se); err != nil {
		return fmt.Errorf("error decoding IDN table ref: %s", err)
	}
	// Write to the output file
	idnRow := []string{idnTableRef.ID, idnTableRef.Url, idnTableRef.UrlPolicy}
	err := a.CSVFiles["idnLanguage"].CsvWriter.Write(StandardizeStringSlice(idnRow))
	if err != nil {
		return err
	}
	return nil
}

// Processes a NNDN object
func processNNDNTag(se *xml.StartElement, a *XMLAnalyzer) error {
	// Found an nndn, add it to the counter for sanity checking
	// Skip nndns that are not in the nndns namespace
	if se.Name.Space != NameSpace["rdeNNDN"] {
		return nil
	}
	a.Counters["nndn"]++
	var nndns XMLNNDN
	if err := a.XMLFile.Decoder.DecodeElement(&nndns, se); err != nil {
		return fmt.Errorf("error decoding nndn: %s", err)
	}
	nndnRow := []string{nndns.AName, nndns.UName, nndns.IDNTableID, nndns.OriginalName, nndns.NameState, nndns.CrDate}
	err := a.CSVFiles["nndn"].CsvWriter.Write(StandardizeStringSlice(nndnRow))
	if err != nil {
		return err
	}
	return nil
}

// Processes a Domain object
func processDomainTag(se *xml.StartElement, a *XMLAnalyzer, uniqueContactIDs map[string]bool) error {

	// Found a domain, add it to the counter for sanity checking
	a.Counters["domain"]++
	// Skip domain tokens that are not in the domain namespace
	if se.Name.Space != NameSpace["rdeDomain"] {
		return nil
	}
	var dom XMLDomain
	if err := a.XMLFile.Decoder.DecodeElement(&dom, se); err != nil {
		return fmt.Errorf("error decoding domain: %s", err)
	}
	// Write the domain to the domain file
	domainRow := []string{string(dom.Name), dom.RoID, dom.UName, dom.IdnTableId, dom.OriginalName, dom.Registrant, dom.ClID, dom.CrRr, dom.CrDate, dom.ExDate, dom.UpRr, dom.UpDate}
	err := a.CSVFiles["domain"].CsvWriter.Write(StandardizeStringSlice(domainRow))
	if err != nil {
		return err
	}
	// Add a line to the contactID file for each contact, only if it does not exist yet
	for _, contact := range dom.Contact {
		// Only add it if it is not there already
		if !uniqueContactIDs[contact.ID] {
			uniqueContactIDs[contact.ID] = true
			a.Counters["uniqueContactID"]++
		}
	}
	// Write the domain statuses to the status file
	dStatuses := []string{dom.Name}
	for _, status := range dom.Status {
		dStatuses = append(dStatuses, status.S)
	}
	for i, s := range dStatuses {
		if i == 0 {
			continue
		}
		a.Counters["domainStatus"]++
		err := a.CSVFiles["domainStatus"].CsvWriter.Write(StandardizeStringSlice([]string{dom.Name, s}))
		if err != nil {
			return err
		}
	}
	// Write the nameservers to the nameserver file
	dNameservers := []string{dom.Name}
	for _, ns := range dom.Ns {
		dNameservers = append(dNameservers, ns.HostObjs...)
	}
	for i, ns := range dNameservers {
		if i == 0 {
			continue
		}
		a.Counters["domainNameservers"]++
		err := a.CSVFiles["domainNameservers"].CsvWriter.Write(StandardizeStringSlice([]string{dom.Name, ns}))
		if err != nil {
			return err
		}
	}
	// Write the dnssec information to the dnssec file
	for _, dsData := range dom.SecDNS.DSData {
		a.Counters["domainDnssec"]++
		dnssecRow := []string{dom.Name, strconv.Itoa(dsData.KeyTag), strconv.Itoa(dsData.Alg), strconv.Itoa(dsData.DigestType), dsData.Digest}
		err := a.CSVFiles["domainDnssec"].CsvWriter.Write(StandardizeStringSlice(dnssecRow))
		if err != nil {
			return err
		}
	}
	// Write the transfer information to the transfer file
	if dom.TrnData.TrStatus.State != "" {
		a.Counters["domainTransfers"]++
		transferRow := []string{dom.Name, dom.TrnData.TrStatus.State, dom.TrnData.ReRr.RegID, dom.TrnData.ReDate, dom.TrnData.ReRr.RegID, dom.TrnData.AcDate, dom.TrnData.ExDate}
		err := a.CSVFiles["domainTransfers"].CsvWriter.Write(StandardizeStringSlice(transferRow))
		if err != nil {
			return err
		}
	}
	return nil
}

// Processes a Host object
func processHostTag(se *xml.StartElement, a *XMLAnalyzer) error {
	// Found a host, add it to the counter for sanity checking
	a.Counters["host"]++
	// Skip host tags that are not in the host namespace
	if se.Name.Space != NameSpace["rdeHost"] {
		return nil
	}
	var host XMLHost
	if err := a.XMLFile.Decoder.DecodeElement(&host, se); err != nil {
		return fmt.Errorf("error decoding host: %s", err)
	}
	hostRow := []string{host.Name, host.RoID, host.ClID, host.CrRr, host.CrDate, host.UpRr, host.UpDate}
	err := a.CSVFiles["host"].CsvWriter.Write(StandardizeStringSlice(hostRow))
	if err != nil {
		return err
	}
	// Set Status in statusFile
	hStatuses := []string{host.Name}
	for _, status := range host.Status {
		a.Counters["hostStatus"]++
		hStatuses = append(hStatuses, status.S)
	}
	for i, s := range hStatuses {
		if i == 0 {
			continue
		}
		err := a.CSVFiles["hostStatus"].CsvWriter.Write(StandardizeStringSlice([]string{host.Name, s}))
		if err != nil {
			return err
		}
	}
	// Set addresses in addrFile
	for _, addr := range host.Addr {
		a.Counters["hostAddress"]++
		err := a.CSVFiles["hostAddress"].CsvWriter.Write(StandardizeStringSlice([]string{host.Name, addr.IP, addr.IP}))
		if err != nil {
			return err
		}
	}
	return nil
}
