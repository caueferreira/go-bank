package db

import (
	"github.com/gocql/gocql"
	"log"
)

func ConnectCassandra() *gocql.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "go_bank"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	return session
}
