# Mini-Scan

This project implements the processing application ("[processor](./cmd/processor/main.go)") for network service scan
data submitted to a GCP Pub/Sub topic by the "[scanner](./cmd/scanner/main.go)" application. The processor pulls scan
results from the corresponding GCP Pub/Sub subscription and maintains an up-to-date record of each unique
`(ip, port, service)` combination and its latest response in CockroachDB. CockroachDB was chosen for its wire
compatibility with PostgreSQL and its distributed, horizontally scalable architecture.

## Implementation Details

This repository was created by forking the original [mini-scan-takehome](https://github.com/censys/mini-scan-takehome)
repo and extending it with additional code for the processor application, database schema migration, end-to-end tests,
and Docker Compose configuration update. The following new files were added to the project:
```
.
├── cmd
│   ├── processor
│   │   ├── Dockerfile
│   │   └── main.go
├── db
│   └── migrations
│       └── V1__initial_schema.sql
├── internal
│   ├── config
│   │   └── processor.go
│   ├── datastores
│   │   ├── datastore_cockroach.go
│   │   ├── datastore_cockroach_test.go
│   │   ├── datastore.go
│   │   ├── datastore_test.go
│   │   └── types.go
│   └── services
│       ├── processor.go
│       └── processor_test.go
└── test
    └── e2e
        ├── e2e_test.go
        ├── features
        │   ├── v1_scan_data.feature
        │   └── v2_scan_data.feature
        ├── go.mod
        └── go.sum
```

The [cmd](./cmd) and [internal](./internal) directories contain the processor application code. The
[cmd/processor/main.go](./cmd/processor/main.go) file is the entry point for the application that handles interactions with
the GCP Pub/Sub subscription. Code for the service and persistence layers is located in the [internal](./internal)
directory.

The [docker-compose.yml](./docker-compose.yml) file was updated to include 3 new services: `cockroach`, `flyway`, and
`processor`. The `cockroach` service runs a single-node CockroachDB instance, the `flyway` service applies the database
schema migration scripts from the [db/migrations](./db/migrations) directory, and the `processor` service runs the
processor application that consumes scan data from the Pub/Sub subscription and stores it in the database.

## Application Configuration

The processor application uses environment variables for its configuration. Although the default values specified in
[internal/config/processor.go](./internal/config/processor.go) are compatible with the provided Docker Compose setup
out of the box, you can override these defaults by setting the appropriate environment variables to adjust the
application's behavior as needed. The following environment variables are supported:

| Environment Variable | Default Value | Description |
| -------------------- | ------------- | ----------- |
| `PROCESSOR_DATASTORE_TYPE` | `cockroachdb` | The type of datastore to use for storing scan data. Currently, only `cockroachdb` is supported. |
| `PROCESSOR_DATASTORE_COCKROACHHOST` | `cockroach` | The hostname of the CockroachDB instance to connect to. |
| `PROCESSOR_DATASTORE_COCKROACHPORT` | `26257` | The port number of the CockroachDB instance to connect to. |
| `PROCESSOR_DATASTORE_COCKROACHUSER` | `root` | The username to use when connecting to the CockroachDB instance. |
| `PROCESSOR_DATASTORE_COCKROACHDATABASE` | `defaultdb` | The name of the database to use for storing scan data. |
| `PROCESSOR_DATASTORE_COCKROACHSCHEMA` | `miniscan` | The schema to use within the database for storing scan data. |
| `PROCESSOR_DATASTORE_COCKROACHTABLE` | `scans` | The name of the table to use within the schema for storing scan data. |
| `PROCESSOR_PUBSUB_PROJECTID` | `mini-scan` | The GCP project ID where the Pub/Sub topic and subscription are located. |
| `PROCESSOR_PUBSUB_SUBSCRIPTIONID` | `scan-results` | The ID of the Pub/Sub subscription that the processor application should pull messages from. |
| `PROCESSOR_PUBSUB_MAXOUTSTANDINGMESSAGES` | `10` | The maximum number of unacknowledged messages that the processor application can receive at once. |
| `PROCESSOR_SERVICE_GRACEFULSHUTDOWNTIMEOUT` | `5s` | The duration to wait for the processor application to gracefully shut down before forcefully terminating it. |

## Testing

To validate behavior of the processor application and ensure data is persisted correctly in the database, it can be
tested either manually or using automated tests. The following sections provide instructions on how to execute both.

### Automated Tests

This project includes unit tests for the components of the processor application, as well as end-to-end tests that cover
the entire data flow from scan data publication to processing and storage in the database.

#### Unit Tests

The current unit test coverage is not perfect, but the implemented tests lay the foundation on which
additional test cases can be added in the future. To run the unit tests, execute the `go test ./...` command in
the repo root directory:
```
$ go test ./...
?       github.com/smolse/scan-takehome/cmd/processor   [no test files]
?       github.com/smolse/scan-takehome/cmd/scanner     [no test files]
?       github.com/smolse/scan-takehome/internal/config [no test files]
?       github.com/smolse/scan-takehome/pkg/scanning    [no test files]
ok      github.com/smolse/scan-takehome/internal/datastores     0.005s
ok      github.com/smolse/scan-takehome/internal/services       0.005s
```

#### End-to-End Tests

End-to-end tests are implemented with the help of the [Godog](https://github.com/cucumber/godog) project, which is an
official [Cucumber](https://cucumber.io/) BDD framework implementation for Golang. It uses
[Gherkin](https://cucumber.io/docs/gherkin/reference/) syntax to define feature files that describe the expected
behavior of the system in a human-readable format. The feature files are located in the
[test/e2e/features](./test/e2e/features) directory. Actual Go code that implements the steps defined in the feature
files is located in the [test/e2e/e2e_test.go](./test/e2e/e2e_test.go) file.

To run the end-to-end tests, execute the following sequence of commands:

1. Ensure all potentially running Docker Compose services are stopped to start the end-to-end tests with a clean state:
   ```
   $ docker-compose down
   ```
2. Start the Docker Compose services except the scanner, as the end-to-end test scenarios will publish scan data to the
Pub/Sub subscription themselves:
    ```
    $ docker-compose up -d pubsub mk-topic mk-subscription cockroach flyway processor
    ```
3. Take a quick look at the Docker Compose logs to ensure all services are running correctly:
    ```
    $ docker-compose logs
    ```
4. Change the current working directory to the `test/e2e` directory:
    ```
    $ cd test/e2e
    ```
5. Run the end-to-end tests using the `go test` command:
    ```
    $ go test
    ```
6. After the tests have completed, stop the Docker Compose services:
    ```
    $ docker-compose down
    ```

### Manual Tests

To manually test the processor application, simply perform these two steps:

1. Start all Docker Compose services, including the scanner:
    ```
    $ docker-compose up -d
    ```
2. Run SQL queries against the CockroachDB instance to inspect the data stored into the database by the processor:
    ```
    $ docker exec -it mini-scan-takehome-cockroach-1 cockroach sql --insecure -e "select * from miniscan.scans"
    ```

## Notes and Ideas

In the interest of time, the current implementation of the processor application is kept simple and isn't completely
production-ready. There are many areas that should be improved in order to achieve true production readiness.

What has already been done properly:
- Graceful shutdown handling was implemented to ensure the processor application can finish processing messages before
  shutting down. Both `SIGINT` (sent by Ctrl+C) and `SIGTERM` (usually sent by container orchestrators) signals are
  caught.

What needs to be improved or added in the future:
- Implement rate limiting for Pub/Sub message processing to have better control over the processing rate and downstream
  resources load.
- Implement exponential backoff and retry logic for data store operations to handle transient errors for improved
  reliability.
- Currently, each scan data record is processed individually. It could be beneficial to batch scan data before storing
  it in the database to reduce the number of database transactions and improve performance.
- Current logging is simplistic. The application needs to be updated to use a structured logging library with log
  levels. One possible option is [slog](https://go.dev/blog/slog), which is a part of the standard library since Go
  version 1.21 (the project is currently using Go 1.20 though). Also, add metrics and traces collection.
- Add a linter, such as [golangci-lint](https://github.com/golangci/golangci-lint), to enforce code quality standards
  and best practices.
