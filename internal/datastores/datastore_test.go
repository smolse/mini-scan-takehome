package datastores

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smolse/scan-takehome/internal/config"
)

// TestNewScanDataStore tests the NewScanDataStore factory function.
func TestNewScanDataStore(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.DataStoreConfig
		wantErr bool
	}{
		{
			name: "CockroachDB data store can be successfully created",
			cfg: &config.DataStoreConfig{
				Type:              "cockroachdb",
				CockroachUser:     "test",
				CockroachHost:     "localhost",
				CockroachPort:     26257,
				CockroachDatabase: "test",
				CockroachSchema:   "test",
				CockroachTable:    "test",
			},
			wantErr: false,
		},
		{
			name: "Unsupported data store type results in an error",
			cfg: &config.DataStoreConfig{
				Type: "unsupported",
			},
			wantErr: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewScanDataStore(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
