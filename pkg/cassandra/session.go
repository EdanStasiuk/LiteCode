package cassandra

import (
	"log"
	"os"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

// Initializes the Cassandra session
// Use CASSANDRA_HOSTS env var (comma-separated) if connecting to a cluster
// Defaults to localhost:9042
func Init() error {
	hostsEnv := os.Getenv("CASSANDRA_HOSTS")
	var hosts []string
	if hostsEnv != "" {
		// Split comma-separated hosts
		hosts = splitAndTrim(hostsEnv)
	} else {
		// Default to localhost for local dev / single-node
		hosts = []string{"127.0.0.1:9042"}
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "litecode"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	Session = session
	log.Println("Cassandra connected to hosts:", hosts)
	return nil
}

// Closes the Cassandra session
func Close() {
	if Session != nil {
		Session.Close()
	}
}

// Splits comma-separated string and trims spaces
func splitAndTrim(s string) []string {
	raw := make([]string, 0)
	for _, h := range split(s, ",") {
		h = trimSpace(h)
		if h != "" {
			raw = append(raw, h)
		}
	}
	return raw
}

// Minimal replacements for strings.Split / strings.TrimSpace to avoid extra imports
func split(s, sep string) []string {
	var res []string
	start := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			res = append(res, s[start:i])
			start = i + len(sep)
		}
	}
	res = append(res, s[start:])
	return res
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
