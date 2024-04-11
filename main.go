package main 

import (
	"fmt"
	"monkey/token"
)


func main () {
	var test token.Token
	test2 := token.Token{Type: "Test", Literal: "again"}
	fmt.Println("Hello there ", test2.Type, test.Literal)
}
