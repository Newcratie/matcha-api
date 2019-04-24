package api

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
	"testing"
)

func TestNull(t *testing.T) {
	fmt.Println("working")
}

func TestGlobal(t *testing.T) {
	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
	app.dbGetTagList()
}
