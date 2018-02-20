package sql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK     // *
	COMMA        // ,
	SINGLEQUOTE  // '
	BACKTICK     // `
	STARTBRACKET // (
	ENDBRACKET   // )
	SEMICOLON    // ;

	// Keywords
	SELECT
	FROM

	INSERT
	INTO
	VALUES
)
