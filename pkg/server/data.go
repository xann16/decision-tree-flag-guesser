package server

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"

	"../flagger"
)

const qdataPath = "data/questions.json"
const ndataPath = "data/dt_nodes.txt"
const traitPath = "data/kb_traits.txt"


func LoadData() {
	loadQuestionData()
	loadDecisionTree()
	loadTags()
}



type QuestionItem struct {
	Text        string
  Notes       []string
}

type QuestionData struct {
	Questions   map[string]QuestionItem
	Notes       map[string]string
	Initials    []string
}

var QData QuestionData

func loadQuestionData() {
	jsonFile, err := os.Open(qdataPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	bs, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(bs, &QData); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
}



var DTree flagger.DecTree

func loadDecisionTree() {
	dtFile, err := os.Open(ndataPath)
  if err != nil {
    panic(err)
  }
  defer dtFile.Close()

	DTree = flagger.LoadDecTree(dtFile)
}



type FlagItem struct {
	Tag   string
	Name  string
}

var FlagList []FlagItem


func loadTags() {
  traitsFile, err := os.Open(traitPath)
  if err != nil {
    panic(err)
  }
  defer traitsFile.Close()

  tm := flagger.LoadTraitMap(traitsFile)

  FlagList = make([]FlagItem, len(tm))

  for tag, node := range tm {
  	FlagList[node.Id] = FlagItem{tag, node.Name}
  }
}
