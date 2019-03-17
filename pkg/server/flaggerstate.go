package server

import (
	"fmt"
	"math/rand"
)

type FlagNode struct {
	Name       string
	Tag        string
	ThumbDir   string
	ItemClass  string
	IsButton   bool
}

type FlaggerState struct {
	Title      string
	Flags      []FlagNode
	Choice     string
	Question   string
	Notes      []string
	IsWalk     bool
	IsFinal    bool
	HasNotes   bool
	YesNode    uint32
	NoNode     uint32
}


func GetStartState() *FlaggerState {
	s := &FlaggerState{}

	s.Title = "Wybór flagi"
	s.Choice = "xx"
	s.Question = "Wybierz jedną z flag..."
	s.Notes = []string{}
	s.IsWalk = false

	s.Flags = make([]FlagNode, len(FlagList))
	for i, n := range FlagList {
		s.Flags[i] = FlagNode{ n.Name, n.Tag, "thumbs", "button-flag", true }
	}

	return s
}

func GetWalkState(choice string, nid uint32) *FlaggerState {
	s := &FlaggerState{}

	level := nidToLevel(nid)

	s.Choice = choice
	s.Notes = []string{}
	s.Flags = make([]FlagNode, len(FlagList))

	rule := DTree.Nodes[nid].Rule

	if rule[0] == 'F' {
		if rule[2:4] == choice {
			s.Question = fmt.Sprintf("Odnaleziono wybraną flagę po %d pytaniach", level - 1)
		} else {
			s.Question = "Poszukiwania zakończone niepowoidzeniem. Upewnij się, że udzielono prawidłowych odpowiedzi na wszystkie pytania."
		}
		s.Title = "Poszukiwanie zakończone"
		s.IsWalk = false
		s.IsFinal = true
	} else {
		qnode := QData.Questions[rule]

		var init string
		if level == 1 {
			init = QData.Initials[0]
		} else {
			init = getRandomInitial()
		}

		s.Question = fmt.Sprintf("PYTANIE #%d: %s %s?", level, init, qnode.Text)
		for i, note := range qnode.Notes {
			s.Notes = append(s.Notes, fmt.Sprintf("[%d] %s", i, QData.Notes[note]))
		}
		s.Title = fmt.Sprintf("Poszukiwanie wybranej flagi - pytanie #%d", level)
		s.YesNode = nid << 1
		s.NoNode = s.YesNode + 1
		s.IsWalk = true
		s.HasNotes = len(s.Notes) > 0
	}

	mask := DTree.Nodes[nid].Mask

	for i, n := range FlagList {
		isChosen := n.Tag == choice
		isActive := mask.Get(uint32(i))

		thDir := "thumbs"
		if !isActive {
			thDir += "_gray"
		}

		var ic string
		if isActive {
			if isChosen {
				ic = "active-chosen-flag"
			} else {
				ic = "active-flag"
			}
		} else {
			if isChosen {
				ic = "rejected-chosen-flag"
			} else {
				ic = "rejected-flag"
			}
		}

		s.Flags[i] = FlagNode{ n.Name, n.Tag, thDir, ic, false }
	}

	return s
}


func nidToLevel(nid uint32) uint32 {
	var res uint32
	for nid != 0 {
		res++
		nid >>= 1
	}
	return res
}

func getRandomInitial() string {
	index := rand.Int() % len(QData.Initials)
	return QData.Initials[index]
}
