package main

import (
	"fmt"
	"os"

	"../pkg/flagger"
)

func main() {
	dtPath := "../data/dt_nodes.txt"

	dtFile, err := os.Open(dtPath)
	if err != nil {
		panic(err)
	}
	defer dtFile.Close()

	dt := flagger.LoadDecTree(dtFile)

	dt.Write(os.Stdout)

	for _, tid := range dt.GetUsedStdTraits() {
		fmt.Println(tid)
	}

}
