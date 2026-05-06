package token

import (
	"fmt"
	"strconv"
)

type Position struct {
	Offset int
	Line   int
	Column int
}

func (pos Position) String() string {
	// this is the function which when called will return
	// file::line::Offset
	s := "main.go"
	if pos.isValid() {
		if s != "" {
			s += ":"
		}
		s += strconv.Itoa(pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	return s
}

func (pos *Position) isValid() bool {
	return pos.Line > 0
}

var lines []int

// this will store the byte offset at which a new line starts
// we will ciopy the logic of addline function here
func SetLinesForContent(input []byte) {
	// how do we refer to our code here ? we dont have a file we are writing to
	lines = []int{0}
	for i, b := range input {
		if b == '\n' {
			lines = append(lines, i+1)
		}
	}
}
