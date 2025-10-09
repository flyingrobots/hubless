package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"

	"github.com/flyingrobots/hubless/internal/cli"
)

func main() {
	root := cli.NewRootCommand()
	if err := fang.Execute(context.Background(), root); err != nil {
		os.Exit(1)
	}
}
