package main

import (
	"fmt"
	"./pkg/server"
)

func main() {
	server.Serve()

	fmt.Println(len(server.QData.Questions))
	fmt.Println(len(server.QData.Notes))
	fmt.Println(len(server.QData.Initials))
	fmt.Println(len(server.DTree.Nodes))
	fmt.Println(len(server.FlagList))
}
