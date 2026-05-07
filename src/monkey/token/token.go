package token

import "strconv"

type TokenType int

type Token struct {
	Type     TokenType
	Literal  string
	Position Position
}

const (
	ILLEGAL TokenType = iota
	EOF
	// identifiers

	IDENT
	INT

	// operators
	ASSIGN
	PLUS

	// DELIMITERS

	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// keywords
	FUNCTION
	LET
	RETURN
	TRUE
	FALSE
	IF
	ELSE
)

var tokenTypes = [...]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	IDENT:     "IDENT",
	INT:       "INT",
	ASSIGN:    "ASSIGN",
	PLUS:      "PLUS",
	COMMA:     "COMMA",
	SEMICOLON: "SEMICOLON",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	FUNCTION:  "FUNCTION",
	LET:       "LET",
	RETURN:    "RETURN",
	TRUE:      "TRUE",
	FALSE:     "FALSE",
	IF:        "IF",
	ELSE:      "ELSE",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok.Type && tok.Type < TokenType(len(tokenTypes)) {
		s = tokenTypes[tok.Type]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok.Type)) + ")"
	}
	return s
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

// we would have used a init function here but we dont need to so we will just focus on the lookup function:
func Lookup(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}
