package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tuxcanfly/mysql2json/sql"
)

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	// Print help
	if path == "" || path == "-h" || path == "--help" {
		log.Printf("Usage: mysql2json <path/to/mysql.dump>")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	stmt, err := sql.NewParser(file).Parse()
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(stmt)
}
