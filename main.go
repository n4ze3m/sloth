package main

import (
	"fmt"
	"os"

	"github.com/nazeemnato/sloth/repl"
)

func main() {
	fmt.Println("Welcome to Sloth")
	repl.Start(os.Stdin, os.Stdout)
}
