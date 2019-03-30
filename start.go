package main

import (
	"fmt"

	"./pkg/server"
)

func main() {
	fmt.Println(1)
	fmt.Println(2)
	fmt.Println(3)

	server.Serve()

	//fmt.Println(len(server.QData.Questions))
	//fmt.Println(len(server.QData.Notes))
	//fmt.Println(len(server.QData.Initials))
	//fmt.Println(len(server.DTree.Nodes))
	//fmt.Println(len(server.FlagList))
}
