package sql

import (
	"fmt"
	"io"
)

// Statement represents a SQL statement.
type Statement interface{}

// SelectStatement represents a SQL SELECT statement.
type SelectStatement struct {
	Fields    []string
	TableName string
}

// InsertStatement represents a SQL INSERT statement.
type InsertStatement struct {
	TableName string
	Values    map[string][]string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (Statement, error) {
	// First token should be a keyword.
	tok, lit := p.scanIgnoreWhitespace()

	switch tok {
	case SELECT:
		return p.parseSelect()
	case INSERT:
		return p.parseInsert()
	default:
		return nil, fmt.Errorf("found %q, expected KEYWORD", lit)
	}
}

func (p *Parser) parseSelect() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	// Return the successfully parsed statement.
	return stmt, nil
}

func (p *Parser) parseInsert() (*InsertStatement, error) {
	stmt := &InsertStatement{
		Values: make(map[string][]string),
	}

	// INTO
	tok, lit := p.scanIgnoreWhitespace()
	if tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}

	// `
	tok, lit = p.scanIgnoreWhitespace()
	if tok != BACKTICK {
		p.unscan()
	}

	// Then we should read the table name.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	// `
	tok, lit = p.scanIgnoreWhitespace()
	if tok != BACKTICK {
		p.unscan()
	}

	// VALUES
	tok, lit = p.scanIgnoreWhitespace()
	if tok != VALUES {
		return nil, fmt.Errorf("found %q, expected VALUES", lit)
	}

	// Next we should loop over all our values of the form
	// ('x', 'y', 'z'), ('a', 'b', 'c'), ....
	for {
		// (
		tok, lit := p.scanIgnoreWhitespace()
		if tok != STARTBRACKET {
			return nil, fmt.Errorf("found %q, expected STARTBRACKET", lit)
		}

		// '
		tok, lit = p.scanIgnoreWhitespace()
		if tok != SINGLEQUOTE {
			return nil, fmt.Errorf("found %q, expected SINGLEQUOTE", lit)
		}

		// Read a value.
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected value", lit)
		}

		key := lit
		var values []string

		// '
		tok, lit = p.scanIgnoreWhitespace()
		if tok != SINGLEQUOTE {
			return nil, fmt.Errorf("found %q, expected SINGLEQUOTE", lit)
		}

		for {
			// If the next token is not a comma then break the loop.
			if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
				p.unscan()
				break
			}

			// '
			tok, lit = p.scanIgnoreWhitespace()
			if tok != SINGLEQUOTE {
				return nil, fmt.Errorf("found %q, expected SINGLEQUOTE", lit)
			}

			// Read a value.
			tok, lit := p.scanIgnoreWhitespace()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected value", lit)
			}

			values = append(values, lit)

			// '
			tok, lit = p.scanIgnoreWhitespace()
			if tok != SINGLEQUOTE {
				return nil, fmt.Errorf("found %q, expected SINGLEQUOTE", lit)
			}
		}

		stmt.Values[key] = values

		// )
		tok, lit = p.scanIgnoreWhitespace()
		if tok != ENDBRACKET {
			return nil, fmt.Errorf("found %q, expected ENDBRACKET", lit)
		}

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA || tok == SEMICOLON {
			p.unscan()
			break
		}
	}

	// Return the successfully parsed statement.
	return stmt, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
