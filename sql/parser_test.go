package sql_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tuxcanfly/mysql2json/sql"
)

// Ensure the parser can parse strings into Select Statement ASTs.
func TestParser_ParseSelectStatement(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *sql.SelectStatement
		err  string
	}{
		// Single field statement
		{
			s: `SELECT name FROM tbl`,
			stmt: &sql.SelectStatement{
				Fields:    []string{"name"},
				TableName: "tbl",
			},
		},

		// Multi-field statement
		{
			s: `SELECT first_name, last_name, age FROM my_table`,
			stmt: &sql.SelectStatement{
				Fields:    []string{"first_name", "last_name", "age"},
				TableName: "my_table",
			},
		},

		// Select all statement
		{
			s: `SELECT * FROM my_table`,
			stmt: &sql.SelectStatement{
				Fields:    []string{"*"},
				TableName: "my_table",
			},
		},

		// Errors
		{s: `foo`, err: `found "foo", expected KEYWORD`},
		{s: `SELECT !`, err: `found "!", expected field`},
		{s: `SELECT field xxx`, err: `found "xxx", expected FROM`},
		{s: `SELECT field FROM *`, err: `found "*", expected table name`},
	}

	for i, tt := range tests {
		stmt, err := sql.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

// Ensure the parser can parse strings into Insert Statement ASTs.
func TestParser_ParseInsertStatement(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *sql.InsertStatement
		err  string
	}{
		// Simple Insert statement
		{
			s: "INSERT INTO `departments` VALUES ('db001', 'Marketing')",
			stmt: &sql.InsertStatement{
				TableName: "departments",
				Values:    map[string][]string{"db001": []string{"Marketing"}},
			},
		},

		// Errors
		{s: `foo`, err: `found "foo", expected KEYWORD`},
	}

	for i, tt := range tests {
		stmt, err := sql.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
