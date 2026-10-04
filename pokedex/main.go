package main

import (
	"os"

	"github.com/yeikol-l/bootdev/pokedex/internal/repl"
)

func main() {
	program := repl.NewProgram(os.Stdin)
	program.Start()
}
