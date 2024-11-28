package services

import (
	"encoding/json"
	"fmt"

	"github.com/smolse/scan-takehome/internal/datastores"
	"github.com/smolse/scan-takehome/pkg/scanning"
)

// ProcessorService is a service that processes scan data and persists it to a data store.
type ProcessorService struct {
	datastore datastores.ScanDataStore
}

// NewProcessorService is a factory function that creates a new ProcessorService instance with the given data store.
func NewProcessorService(dataStore datastores.ScanDataStore) *ProcessorService {
	return &ProcessorService{
		datastore: dataStore,
	}
}

// transformScanData transforms the scan data from the scanning package to the data store format.
func (s *ProcessorService) transformScanData(scan *scanning.Scan) (datastores.Scan, error) {
	var response string
	dataBytes, err := json.Marshal(scan.Data)
	if err != nil {
		return datastores.Scan{}, fmt.Errorf("failed to marshal scan data: %w", err)
	}

	switch scan.DataVersion {
	case scanning.V1:
		var v1Data scanning.V1Data
		err = json.Unmarshal(dataBytes, &v1Data)
		if err != nil {
			return datastores.Scan{}, fmt.Errorf("failed to parse V1 data: %w", err)
		}
		response = string(v1Data.ResponseBytesUtf8)
	case scanning.V2:
		var v2Data scanning.V2Data
		err = json.Unmarshal(dataBytes, &v2Data)
		if err != nil {
			return datastores.Scan{}, fmt.Errorf("failed to parse V2 data: %w", err)
		}
		response = v2Data.ResponseStr
	default:
		return datastores.Scan{}, fmt.Errorf("unsupported data version: %d", scan.DataVersion)
	}

	return datastores.Scan{
		Ip:        scan.Ip,
		Port:      scan.Port,
		Service:   scan.Service,
		Timestamp: scan.Timestamp,
		Response:  response,
	}, nil
}

// ProcessScanData processes the consumed scan data and persists it to the data store.
func (s *ProcessorService) ProcessScanData(scanData []byte) error {
	// Unmarshal the raw scan data
	scan := &scanning.Scan{}
	err := json.Unmarshal(scanData, scan)
	if err != nil {
		return fmt.Errorf("failed to unmarshal scan data: %w", err)
	}

	// Transform the scan data to the data store format
	scanRecord, err := s.transformScanData(scan)
	if err != nil {
		return fmt.Errorf("failed to transform scan data to the data store representation: %w", err)
	}

	// TODO: Add exponential backoff retry logic for transient errors
	err = s.datastore.UpdateScanData(scanRecord)
	if err != nil {
		return fmt.Errorf("failed to update scan data: %w", err)
	}
	return nil
}
