package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Context struct{}

type Globals struct {
	Version VersionFlag `name:"version" help:"Print version information."`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Printf("v%s\n", vars["version"])
	app.Exit(0)
	return nil
}

type InitCmd struct {
	Verbose int `short:"v" type:"counter" help:"Verbosity can have a value of 1-3. Example: --verbose=3 or -vvv."`
}

func (cmd *InitCmd) Validate() error {
	return nil
}

func (cmd *InitCmd) Run(ctx *Context) error {
	return nil
}

type CLI struct {
	Globals
	Init InitCmd `cmd:""`
}
