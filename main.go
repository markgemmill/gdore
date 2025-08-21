package main

import (
	"fmt"
	"gdore/environ"
	"gdore/gui"
)

func main() {
	env, err := environ.CreateEnvironment()
	if err != nil {
		fmt.Printf("error loading environment settings: %s\n", err)
		return
	}

	gui.Gui(env)
}
