Feature: V1 scan data

  Background:
    Given PubSub client for project "test-project" and topic "scan-topic" is ready
    And the Postgres client for "postgresql://root@localhost:26257/defaultdb" database is ready

  Scenario: Valid V1 scan data for a new service is stored in the database

    Given the "miniscan.scans" table is empty
    When a scan message with the following payload is published to the "scan-topic" topic:
      """
      {
        "ip": "2.2.2.2",
        "port": 80,
        "service": "HTTP",
        "timestamp": 1732769424,
        "data_version": 1,
        "data": {
          "response_bytes_utf8": "aGVsbG8gd29ybGQ="
        }
      }
      """
    Then the following data can be found in the "miniscan.scans" table after 3 seconds:
      | ip       | port | service | timestamp   | response    |
      | 2.2.2.2  | 80   | HTTP    | 1732769424  | hello world |
