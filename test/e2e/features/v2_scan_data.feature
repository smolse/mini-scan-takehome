Feature: V2 scan data

  Background:
    Given PubSub client for project "test-project" and topic "scan-topic" is ready
    And the Postgres client for "postgresql://root@localhost:26257/defaultdb" database is ready

  Scenario: Valid V2 scan data for a new service is stored in the database

    Given the "miniscan.scans" table is empty
    When a scan message with the following payload is published to the "scan-topic" topic:
      """
      {
        "ip": "3.3.3.3",
        "port": 21,
        "service": "FTP",
        "timestamp": 1732774035,
        "data_version": 2,
        "data": {
          "response_str": "hello world"
        }
      }
      """
    Then the following data can be found in the "miniscan.scans" table after 3 seconds:
      | ip       | port | service | timestamp   | response    |
      | 3.3.3.3  | 21   | FTP     | 1732774035  | hello world |
