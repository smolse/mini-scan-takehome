package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smolse/scan-takehome/internal/datastores"
)

// MockScanDataStore is a mock implementation of the ScanDataStore interface.
type MockScanDataStore struct {
	mock.Mock
}

// Connect is a mock implementation of the Connect method.
func (m *MockScanDataStore) Connect() error {
	args := m.Called()
	return args.Error(0)
}

// Close is a mock implementation of the Close method.
func (m *MockScanDataStore) Close() error {
	args := m.Called()
	return args.Error(0)
}

// UpdateScanData is a mock implementation of the UpdateScanData method.
func (m *MockScanDataStore) UpdateScanData(scanData datastores.Scan) error {
	args := m.Called(scanData)
	return args.Error(0)
}

// TestProcessScanData tests the ProcessScanData method of the ProcessorService.
func TestProcessScanData(t *testing.T) {
	mockDataStore := new(MockScanDataStore)
	processorService := NewProcessorService(mockDataStore)

	tests := []struct {
		name     string
		scanData []byte
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid scan data V1 is processed successfully",
			scanData: []byte(`{"ip":"1.1.1.1","port":80,"service":"HTTP","timestamp":0,"data_version": 1,"data": {"response_bytes_utf8": "aGVsbG8gd29ybGQ="}}`),
			wantErr:  false,
		},
		{
			name:     "Valid scan data V2 is processed successfully",
			scanData: []byte(`{"ip":"1.1.1.1","port":80,"service":"HTTP","timestamp":0,"data_version": 2,"data": {"response_str": "hello world"}}`),
			wantErr:  false,
		},
		{
			name:     "Invalid scan data results in an error",
			scanData: []byte(`invalid`),
			wantErr:  true,
			errMsg:   "failed to unmarshal scan data",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockDataStore.On("UpdateScanData", mock.Anything).Return(nil).Once()
			}

			err := processorService.ProcessScanData(tt.scanData)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				mockDataStore.AssertCalled(t, "UpdateScanData", mock.Anything)
			}

			mockDataStore.AssertExpectations(t)
		})
	}
}
