package flagger

import (
  "os"
  "fmt"
  "sort"
  "bufio"
  "strings"
  "strconv"
)

type DTNode struct {
  Rule   string
  Mask   Bits256
}

type DecTree struct {
  Levels uint32
  Nodes  []DTNode
}



func ( dt DecTree ) Write(fout *os.File) {
  fmt.Fprintf(fout, "%d\n", dt.Levels)
  for i, n := range dt.Nodes {
    if i != 0 {
      //fmt.Fprintf(fout, "%4d %s %s\n", i, n.mask.ToString(), n.rule)
      fmt.Fprintf(fout, "%4d %s %8s %4d\n", i, n.Mask.ToString(), n.Rule, n.Mask.Count())
    }
  }
}









func ( dt DecTree ) GetUsedStdTraits() []int {
  set := map[int]int{}
  for _, n := range dt.Nodes {
    rule := n.Rule
    if len(rule) > 0 && rule[0] == 'S' {
      val, err := strconv.ParseInt(rule[2:], 10, 32)
      if err != nil {
        panic(err)
      }
      set[int(val)] = 1
    }
  }
  res := []int{}
  for k, _ := range set {
    res = append(res, k)
  }
  sort.Ints(res)
  return res
}








func GetEmptyDecTree( leafCount, spareLevels uint32 ) DecTree {
  levels := calculateLevels(leafCount, spareLevels)
  Nodes := int(1 << levels)

  return DecTree{ levels, make([]DTNode, Nodes, Nodes) }
}

func LoadDecTree( fin *os.File ) DecTree {
  scanner := bufio.NewScanner(fin)
  if !scanner.Scan() {
    panic(fmt.Errorf("No header!"))
  }
  header := scanner.Text()
  levelCount, err := strconv.ParseUint(header, 10, 32)
  if err != nil {
    panic(err)
  }

  levels := uint32(levelCount)
  Nodes := int(1 << levels)

  dt := DecTree{ levels, make([]DTNode, Nodes, Nodes) }

  for scanner.Scan() {
    line := scanner.Text()
    if len(line) == 0 || line[0] == '#' {
      continue
    }

    fields := strings.Fields(line)

    index, err := strconv.ParseUint(fields[0], 10, 32)
    if err != nil {
      panic(err)
    }
    mask := GetBits256FromString(fields[1])
    rule := fields[2]

    dt.Nodes[index] = DTNode{ rule, mask }
  }

  return dt
}









func calculateLevels( leafCount, spareLevels uint32 ) uint32 {
  var pow uint32
  val := uint32(1)
  for val < leafCount {
    val <<= 1
    pow++
  }
  return pow + spareLevels + 1
}
