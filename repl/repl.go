package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/parser"
)

const PROMT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out,"Woops!! We ran into some errors\n")
		io.WriteString(out,"parser errors:\n")
		io.WriteString(out, "\t"+msg+"\n")
	}
}
