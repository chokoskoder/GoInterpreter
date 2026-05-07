package lexer

import (
	"testing"
	"monkey/token"
)

func TestNextToken(t *testing.T){

	input := "=+{}(),;"

	tests := []struct{
		expectedType 		token.TokenType
		expectedLiteral	string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
}
}

for i, tt := range tests {
	tok := l.NexToken()

	if tok.Type
}
