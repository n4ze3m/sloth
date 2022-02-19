package main

import (
	// "strings"
	"syscall/js"

	"github.com/nazeemnato/sloth/evaluator"
	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/object"
	"github.com/nazeemnato/sloth/parser"
)

func startWasm(input string) string {
	env := object.NewEnviroment()
	// split input by \n
	// lines := strings.Split(input, "\n")
	// var output string
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return parseErrors(p.Errors())
	}
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		return evaluated.Inspect() + "\n"
	}
	return ""
	// loop through each line
	// for _, line := range lines {
	// 	l := lexer.New(line)
	// 	p := parser.New(l)
	// 	program := p.ParseProgram()
	// 	if len(p.Errors()) != 0 {
	// 		return "parser errors"
	// 	}
	// 	evaluated := evaluator.Eval(program, env)

	// 	if evaluated != nil {
	// 		output += evaluated.Inspect() + "\n"
	// 	}
	// }
	// return output
}

func wasmWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := args[0].String()
		output := startWasm(input)
		return output
	})
}

func main() {
	js.Global().Set("mySloth", wasmWrapper())
	<-make(chan bool)
}

func parseErrors(errors []string) string {
	var err string
	for _, msg := range errors {
		err += "parser errors:\n"
		err += "\t" + msg + "\n"
	}
	return err
}
