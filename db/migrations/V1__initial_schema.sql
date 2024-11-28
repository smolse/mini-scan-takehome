CREATE SCHEMA IF NOT EXISTS miniscan;

CREATE TABLE IF NOT EXISTS miniscan.scans (
    ip INET,
    port INT,
    service STRING,
    timestamp TIMESTAMP,
    response STRING,
    PRIMARY KEY (ip, port, service)
);
