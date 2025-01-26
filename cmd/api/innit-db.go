package main

import (
	"github.com/gocql/gocql"
	"log"
)

func InnitDatabase() {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer session.Close()

	if err := initializeDatabase(session); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
}

func initializeDatabase(session *gocql.Session) error {
	cqlCommands := []string{
		`CREATE KEYSPACE IF NOT EXISTS go_bank WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};`,
		`DROP TABLE IF EXISTS go_bank.accounts`,
		`DROP TABLE IF EXISTS go_bank.transactions`,
		`DROP TABLE IF EXISTS go_bank.transfers`,
		`CREATE TABLE IF NOT EXISTS go_bank.accounts (id UUID PRIMARY KEY, name TEXT, email TEXT, sort_code TEXT, account_number TEXT, balance INT);`,
		`CREATE TABLE IF NOT EXISTS go_bank.transactions (id UUID PRIMARY KEY, account_id TEXT, transaction_tyoe TEXT, amount INT);`,
		`CREATE TABLE IF NOT EXISTS go_bank.transfers (id UUID PRIMARY KEY, to_account UUID, from_account UUID, amount INT, success BOOLEAN);`,
	}

	for _, cmd := range cqlCommands {
		if err := session.Query(cmd).Exec(); err != nil {
			return err
		}
	}
	return nil
}
