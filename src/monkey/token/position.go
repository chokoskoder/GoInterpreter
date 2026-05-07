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

func InitLines() {
	lines = []int{0}
}

func AddLine(offset int) {
	lines = append(lines, offset)
}

func CalculatePosition(offset int) (Position, error) {
	line := -1
	column := 0

	// Binary search our way through
	high := len(lines) - 1
	low := 0
	for low <= high {
		mid := (high + low) / 2
		if lines[mid] == offset {
			// this is an unlikely case , but still we can exit here and simply ut the value in line
			line = mid + 1
			break
		} else if lines[mid] > offset {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	// now this is where we need to check where is our high and low pointer and decide  where our offset lies
	// the loop will only finish once high is no longer higher than low and that will mean that high is now low and ideally
	// that is where our offset should lie
	if line == -1 {
		line = high + 1
	}

	// now we need to start calculating the column:
	// it is going to be the difference in our lines initial offset and the offset we have been given
	column = offset - lines[line-1]

	return Position{
		offset,
		line,
		column,
	}, nil
}
