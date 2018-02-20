package main

import (
	"fmt"
	"strings"

	"github.com/tuxcanfly/mysql2json/sql"
)

func main() {
	str := `SELECT * FROM my_table`
	stmt, err := sql.NewParser(strings.NewReader(str)).Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(stmt)
}
