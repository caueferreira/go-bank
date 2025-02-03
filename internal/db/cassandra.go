package db

import (
	"github.com/gocql/gocql"
	"log"
	"sync"
)

var session *gocql.Session
var once sync.Once

func GetSession() *gocql.Session {
	once.Do(func() {
		cluster := gocql.NewCluster("localhost")
		cluster.Port = 9042
		cluster.Keyspace = "go_bank"
		cluster.Consistency = gocql.Quorum

		var err error
		session, err = cluster.CreateSession()
		if err != nil {
			log.Fatal("Failed to create session:", err)
		}
	})

	return session
}

func CloseCassandraSession() {
	session.Close()
}
