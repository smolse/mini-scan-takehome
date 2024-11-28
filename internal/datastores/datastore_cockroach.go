package datastores

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/smolse/scan-takehome/internal/config"
)

// PgxConnectionPool is an interface for a pgx connection pool.
type PgxConnectionPool interface {
	Close()
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

// CockroachScanDataStore is a data store implementation that persists scan data to a CockroachDB cluster.
type CockroachScanDataStore struct {
	database string
	schema   string
	table    string

	pool PgxConnectionPool
}

// NewCockroachScanDataStore creates a new CockroachScanDataStore instance based on the provided configuration.
func NewCockroachScanDataStore(cfg *config.DataStoreConfig) (*CockroachScanDataStore, error) {
	return &CockroachScanDataStore{
		database: fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			cfg.CockroachUser,
			cfg.CockroachHost,
			strconv.Itoa(cfg.CockroachPort),
			cfg.CockroachDatabase,
		),
		schema: cfg.CockroachSchema,
		table:  cfg.CockroachTable,
	}, nil
}

// Connect establishes a connection to the CockroachDB cluster.
func (s *CockroachScanDataStore) Connect() error {
	pool, err := pgxpool.Connect(context.Background(), s.database)
	if err != nil {
		return fmt.Errorf("failed to connect to CockroachDB: %w", err)
	}
	s.pool = pool
	return nil
}

// Close closes the connection pool for the CockroachDB cluster.
func (s *CockroachScanDataStore) Close() error {
	s.pool.Close()
	return nil
}

// UpdateScanData updates the scan data in the CockroachDB cluster.
func (s *CockroachScanDataStore) UpdateScanData(scanData Scan) error {
	_, err := s.pool.Exec(context.Background(),
		fmt.Sprintf(`
			UPSERT INTO %s.%s (ip, port, service, timestamp, response)
			VALUES ($1, $2, $3, to_timestamp($4), $5)
		`, s.schema, s.table),
		scanData.Ip,
		scanData.Port,
		scanData.Service,
		scanData.Timestamp,
		scanData.Response,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert scan data: %w", err)
	}
	return nil
}
