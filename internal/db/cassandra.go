package db

import (
	"github.com/gocql/gocql"
	"log"
)

var Session *gocql.Session

func ConnectCassandra() {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "go_bank"
	cluster.Consistency = gocql.Quorum
	var err error
	
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
}

func CloseCassandraSession() {
	Session.Close()
}
