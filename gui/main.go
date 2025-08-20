package main

import (
	"fmt"
	"gdore/environ"
)

func main() {
	env, err := environ.CreateEnvironment()
	if err != nil {
		fmt.Printf("error loading environment settings: %s\n", err)
		return
	}

	Gui(env)
}
