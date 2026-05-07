package lexer

import (
	"github.com/chokoskoder/GoInterpreter/token"
)

type Lexer struct {
	input    string
	position token.Position
}
