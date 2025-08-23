package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli,
		kong.Name("gdore"),
		kong.Description("Guillaume's App"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"version": "0.4.0",
		})

	err := ctx.Run(&Context{})
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
