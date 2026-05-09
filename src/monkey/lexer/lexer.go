package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/chokoskoder/GoInterpreter/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune

	ErrorCount int
	err        func(pos token.Position, msg string)
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.next()
	token.InitLines()
	return l
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // end of file
)

// how does this function return a byte when all the chars are been represented in runes ?
func (l *Lexer) peek() byte {
	if l.readPosition < len(l.input) {
		return l.input[l.readPosition]
	}
	return 0
}

func (l *Lexer) next() {
	if l.readPosition < len(l.input) {
		l.position = l.readPosition
		if l.ch == '\n' {
			// TODO: track line and column position here
			token.AddLine(l.position)
		}

		r, w := rune(l.input[l.readPosition]), 1
		switch {
		case r == 0:
			l.error(l.position, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// character is > 127 bytes, so it's a multi-byte unicode char
			r, w = utf8.DecodeRuneInString(l.input[l.readPosition:])
			if r == utf8.RuneError && w == 1 {
				in := l.input[l.readPosition:]
				if l.position == 0 &&
					len(in) >= 2 &&
					(in[0] == 0xFF && in[1] == 0xFE || in[0] == 0xFE && in[1] == 0xFF) {
					// U+FEFF BOM at start of file, encoded as big- or little-endian
					// UCS-2 (i.e. 2-byte UTF-16). Give specific error (go.dev/issue/71950).
					l.error(l.position, "illegal UTF-8 encoding (got UTF-16)")
					l.readPosition += len(in) // consume all input to avoid error cascade
				} else {
					l.error(l.position, "Illegal utf-8 encoding")
				}
			} else if r == bom && l.position > 0 {
				l.error(l.position, "illegal byte order mark")
			}
		}
		l.readPosition += w
		l.ch = r

	} else {
		l.position = len(l.input)
		if l.ch == '\n' {
			// again here we need to check the logic of how we are going to handle new lines and ensure that
			// we are storing each lines initial byte -> POSITION MARKING
			// TODO -> implement a way to add new line marking here
			token.AddLine(l.position)
		}
		l.ch = eof
	}
}

func (l *Lexer) error(position int, msg string) {
	if l.err != nil {
		l.err(token.Position{Offset: position}, msg)
	}
	l.ErrorCount++
}

func (l *Lexer) errorf(position int, format string, arg ...any) {
	l.error(position, fmt.Sprintf(format, arg...))
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.next()
	}
}

func (l *Lexer) ScanComments() (string, int) {
	offs := l.position - 1
	next := -1
	numCR := 0 // this is for \r
	nlOffset := 0

	if l.ch == '/' {
		//=style comment
		// (the final '\n' is not considered part of the comment)

		l.next()

		for l.ch != '\n' && l.ch >= 0 {
			if l.ch == '\r' {
				numCR++
			}
			l.next()
		}
	}
}

// the library is using insertSemi for some reason, I have not figured it out yet and I need to
func (l *Lexer) Scan() (returnToken token.Token) {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case '=':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.ASSIGN, l.ch, pos)
		}
	case '+':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.PLUS, l.ch, pos)
		}
	case '{':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.LBRACE, l.ch, pos)
		}
	case '}':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.RBRACE, l.ch, pos)
		}
	case '(':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.LPAREN, l.ch, pos)
		}
	case ')':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.RPAREN, l.ch, pos)
		}
	case ',':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.COMMA, l.ch, pos)
		}
	case ';':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.SEMICOLON, l.ch, pos)
		}
	case '!':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.BANG, l.ch, pos)
		}
	case '-':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.MINUS, l.ch, pos)
		}
	case '/':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.SLASH, l.ch, pos)
		}
	case '*':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.ASTERISK, l.ch, pos)
		}
	case '<':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.LT, l.ch, pos)
		}
	case '>':
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = newToken(token.GT, l.ch, pos)
		}
	case eof:
		pos, err := token.CalculatePosition(l.position)
		if err != nil {
			l.errorf(l.position, "%s", err.Error())
		} else {
			tok = token.Token{token.EOF, string(""), pos}
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.Lookup(tok.Literal)
			pos, err := token.CalculatePosition(l.position)
			if err == nil {
				tok.Position = pos
			} else {
				l.errorf(l.position, "%s", err.Error())
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			pos, err := token.CalculatePosition(l.position)
			if err == nil {
				tok.Position = pos
			} else {
				l.errorf(l.position, "%s", err.Error())
			}
			return tok
		} else {
			pos, err := token.CalculatePosition(l.position)
			if err == nil {
				tok.Position = pos
			} else {
				l.errorf(l.position, "%s", err.Error())
			}
			tok = newToken(token.ILLEGAL, l.ch, pos)
		}

	}
	l.next()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.next()
	}
	return l.input[position:l.position]
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func lower(ch rune) rune {
	return ('a' - 'A' | ch)
}

// this is really interesting and uses mathematics, nice dsa solulu

func newToken(tokenType token.TokenType, ch rune, charPos token.Position) token.Token {
	return token.Token{tokenType, string(ch), charPos}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.next()
	}
	return l.input[position:l.position]
}
