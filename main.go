package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello there %s, home dir: %s\n", user.Username, user.HomeDir)

	repl.Start(os.Stdin, os.Stdout)

}
