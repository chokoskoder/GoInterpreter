package repl

import (
	"bufio"
	"fmt"
	"io"
	"github.com/chokoskoder/GoInterpreter/lexer"
	"github.com/chokoskoder/GoInterpreter/token"
)



const PROMPT = ">>"

func Start(in io.Reader, out io.Writer){
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewWithMode(line, lexer.ScanComments)
		for tok := l.Scan(); tok.Type != token.EOF; tok = l.Scan() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
