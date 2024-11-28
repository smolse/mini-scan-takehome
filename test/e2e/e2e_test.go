package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v4"
)

var (
	pubsubClient *pubsub.Client
	db           *pgx.Conn
)

func initializePubSubClient(projectID, topicID string) error {
	if err := os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085"); err != nil {
		return err
	}

	var err error
	pubsubClient, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return err
	}
	return nil
}

func initializePostgresClient(connStr string) error {
	var err error
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	return nil
}

func publishMessage(topicID string, payload string) error {
	topic := pubsubClient.Topic(topicID)
	result := topic.Publish(context.Background(), &pubsub.Message{
		Data: []byte(payload),
	})
	_, err := result.Get(context.Background())
	return err
}

func tableIsEmpty(tableName string) error {
	_, err := db.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", tableName))
	return err
}

func dataCanBeFoundInTable(tableName string, ip string, port int, service string, timestamp int64, response string) error {
	var count int
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(`
			SELECT COUNT(*)
			FROM %s
			WHERE ip = $1 AND port = $2 AND service = $3 AND timestamp = to_timestamp($4) AND response = $5
		`, tableName), ip, port, service, timestamp, response,
	).Scan(&count)

	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("expected 1 row, got %d", count)
	}
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^the "([^"]*)" table is empty$`, tableIsEmpty)
	sc.Step(`^PubSub client for project "([^"]*)" and topic "([^"]*)" is ready$`, initializePubSubClient)
	sc.Step(`^the Postgres client for "([^"]*)" database is ready$`, initializePostgresClient)
	sc.Step(`^a scan message with the following payload is published to the "([^"]*)" topic:$`, publishMessage)
	sc.Step(`^the following data can be found in the "([^"]*)" table after (\d+) seconds:$`, func(tableName string, seconds int, table *godog.Table) error {
		time.Sleep(time.Duration(seconds) * time.Second)
		for _, row := range table.Rows[1:] {
			ip := row.Cells[0].Value
			fmt.Println(row.Cells[1].Value)
			port, err := strconv.Atoi(row.Cells[1].Value)
			if err != nil {
				return err
			}
			service := row.Cells[2].Value
			timestamp, err := strconv.ParseInt(row.Cells[3].Value, 10, 64)
			if err != nil {
				return err
			}
			response := row.Cells[4].Value
			err = dataCanBeFoundInTable(tableName, ip, port, service, timestamp, response)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
