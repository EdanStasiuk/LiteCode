package cassandra

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

// Init initializes a Cassandra session, retrying a few times if needed.
func Init() error {
	hostsEnv := os.Getenv("CASSANDRA_HOSTS")
	var hosts []string
	if hostsEnv != "" {
		hosts = strings.Split(hostsEnv, ",")
	} else {
		hosts = []string{"cassandra1", "cassandra2", "cassandra3"}
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "litecode"
	cluster.Consistency = gocql.Quorum

	var err error
	for i := 0; i < 10; i++ {
		Session, err = cluster.CreateSession()
		if err == nil {
			log.Println("Cassandra connected to hosts:", hosts)
			return nil
		}
		log.Println("Cassandra not ready, retrying in 3s:", err)
		time.Sleep(3 * time.Second)
	}

	return err
}

// Close closes the Cassandra session
func Close() {
	if Session != nil {
		Session.Close()
	}
}
