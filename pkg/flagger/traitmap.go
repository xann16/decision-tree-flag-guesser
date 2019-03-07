package flagger

import (
  "sort"
  "os"
  "fmt"
  "bufio"
  "strings"
  "strconv"
)

type TraitNode struct {
  Id     uint32
  Traits uint64
  Name   string
}

type TraitMap map[string]TraitNode


func ( tm TraitMap ) Update( tag, name string, traits uint64 ) {
  if n, ok := tm[tag]; ok {
    n.Name = name
    n.Traits = traits
    } else {
      panic(fmt.Errorf("No traits for tag %q.", tag))
    }
}

func ( tm TraitMap ) Add( tag, name string, id uint32, traits uint64 ) {
  tm[tag] = TraitNode{ id, traits, name }
}

func ( tm TraitMap ) WriteOrdered( fout *os.File, separated bool ) {
  tags := []string{}
  for tag, _ := range tm {
    tags = append(tags, tag)
  }
  sort.Strings(tags)
  for id, tag := range tags {
    node := tm[tag]
    fmt.Fprintf(fout, "%3d  %2s   %s %s\n",
                id, tag, traitsToString(node.Traits, separated), node.Name)
  }
}

func traitsToString( traits uint64, separated bool ) string {
  var res string
  for i := uint(0); i < 64; i++ {
    if separated && ( ( i % 8 ) == 0 ) {
      res += "#"
    }
    if ( traits & ( 1 << i ) ) == 0 {
      res += "0"
    } else {
      res += "1"
    }
  }
  res += "#"
  return res
}

func stringToTraits( str string ) ( traits uint64 ) {
  var pos uint32
  for _, rune := range str {
    switch rune {
      case '#':
        continue
      case '1':
        traits += ( 1 << pos )
        pos++
      case '0':
        pos++
      default:
        panic(fmt.Errorf("Unrecognized character in traits string"))
    }
  }
  if pos != 64 {
    panic(fmt.Errorf("Invalid traits string length (%d).", pos))
  }
  return
}

func CreateTraitMap( tags []string ) TraitMap {
  tm := TraitMap{}
  for _, tag := range tags {
    tm[tag] = TraitNode{}
  }
  return tm
}

func LoadTraitMap( fin *os.File ) TraitMap {
  tm := TraitMap{}
  scanner := bufio.NewScanner(fin)

  for scanner.Scan() {
    line := scanner.Text()
    if len(line) == 0 || line[0] == '#' {
      continue
    }

    fields := strings.Fields(line)

    id, err := strconv.ParseUint(fields[0], 10, 32)
    if err != nil {
      panic(err)
    }
    tag := fields[1]
    name := strings.Join(fields[3:], " ")
    traits := stringToTraits(fields[2])
    tm.Add(tag, name, uint32(id), traits)
  }

  return tm
}
