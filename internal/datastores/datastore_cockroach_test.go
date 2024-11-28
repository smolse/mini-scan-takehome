package datastores

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPgxConnectionPool is a mock implementation of the PgxConnectionPool interface.
type MockPgxConnectionPool struct {
	mock.Mock
}

// Close is a mock implementation of the Close method.
func (m *MockPgxConnectionPool) Close() {
	m.Called()
}

// Exec is a mock implementation of the Exec method.
func (m *MockPgxConnectionPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

// TestUpdateScanData tests the UpdateScanData method of the CockroachScanDataStore.
func TestUpdateScanData(t *testing.T) {
	mockPgxConnPool := new(MockPgxConnectionPool)
	dataStore := CockroachScanDataStore{
		database: "test",
		schema:   "test",
		table:    "test",
		pool:     mockPgxConnPool,
	}

	tests := []struct {
		name    string
		scan    Scan
		wantErr bool
	}{
		{
			name: "Valid scan data is updated successfully",
			scan: Scan{
				Ip:        "1.1.1.1",
				Port:      80,
				Service:   "HTTP",
				Timestamp: 0,
				Response:  "hello world",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockPgxConnPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil).Once()
			}

			err := dataStore.UpdateScanData(tt.scan)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
