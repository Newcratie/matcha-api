package kwal

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"os"
	"time"
)

type Key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func newConn(host string) (bolt.Conn, error) {
	DB, err := bolt.NewDriverPool("bolt://neo4j:secret@"+host+":7687", 1000)
	if err != nil {
		return nil, err
	}
	retries := 0
	for retries < 300 {
		conn, _ := DB.OpenPool()
		if conn != nil {
			return conn, nil
		}
		retries = retries + 1
		fmt.Println("KWAL neo4j not ready, waiting 1s and trying again ", retries)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

func GetKeys() (Keys []Key) {
	host := os.Getenv("NEO_HOST")
	conn, _ := newConn(host)
	data, _, _, _ := conn.QueryNeoAll(`MATCH (n:Kwal) RETURN n`, nil)
	for _, node := range data {
		Keys = append(Keys, Key{
			node[0].(graph.Node).Properties["name"].(string),
			node[0].(graph.Node).Properties["value"].(string),
		})
	}
	fmt.Println("Keys: ", Keys)
	return
}
