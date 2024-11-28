package datastores

// Scan is a struct type that represents a scan record in the data store.
type Scan struct {
	Ip        string
	Port      uint32
	Service   string
	Timestamp int64
	Response  string
}
