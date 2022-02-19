package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nazeemnato/sloth/evaluator"
	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/object"
	"github.com/nazeemnato/sloth/parser"
)

const PROMT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnviroment()

	for {
		fmt.Print(PROMT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == ".exit" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluted := evaluator.Eval(program, env)
		
		if evaluted != nil {
			io.WriteString(out, evaluted.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out,"Woops!! We ran into some errors\n")
		io.WriteString(out,"parser errors:\n")
		io.WriteString(out, "\t"+msg+"\n")
	}
}
