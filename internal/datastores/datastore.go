package datastores

import (
	"fmt"

	"github.com/smolse/scan-takehome/internal/config"
)

// ScanDataStore is an interface for data stores where scan data is persisted after processing.
type ScanDataStore interface {
	Connect() error
	Close() error

	UpdateScanData(Scan) error
}

// NewScanDataStore is a factory function that creates a new ScanDataStore instance based on the provided configuration.
func NewScanDataStore(cfg *config.DataStoreConfig) (ScanDataStore, error) {
	switch cfg.Type {
	case "cockroachdb":
		return NewCockroachScanDataStore(cfg)
	default:
		return nil, fmt.Errorf("unsupported data store type: %s", cfg.Type)
	}
}
