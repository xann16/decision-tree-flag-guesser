package flagger

import (
  "fmt"
  "math"
  //"strings"
)

func BuildDecTree( leafCount, spareLevels uint32,
                   traits *TraitMatrix,
                   tm TraitMap,
                   splitPercTolerance float32) DecTree {
  dt := GetEmptyDecTree(leafCount, spareLevels)

  tags := make([]string, len(tm))
  for tag, n := range tm {
    tags[n.Id] = tag
  }

  calculateNode( 1, 1, leafCount, GetOnes256(leafCount), traits, &dt, tags )

  return dt
}

func calculateNode( level, index, count uint32, mask Bits256,
                    traits *TraitMatrix, dt *DecTree, tags []string ) {
  //fmt.Println("LEVEL:", level, "NODE:", index, "COUNT:", count)

  if count == 1 {
    createLeaf(index, mask, dt, tags)
    return
  }

  if level == dt.Levels {
    panic(fmt.Errorf("Reached the bottom"))
  }

  newMask, traitId, newCount, quality := findBestStdFit(count, mask, traits)

  if quality > ((float64(count) / 2) - 0.25) {
    createDeadlock(index, count, mask, dt, tags)
    return
  }

  createStandard(index, traitId, mask, dt)

  calculateNode(level + 1, index << 1, newCount, newMask, traits, dt, tags)
  calculateNode(level + 1, (index << 1) + 1, count - newCount,
                And(Not(newMask), mask), traits, dt, tags)
}


func findBestStdFit(count uint32,
                    mask Bits256,
                    traits *TraitMatrix ) ( Bits256, uint32, uint32, float64 ) {
  var bestMask Bits256
  var bestIndex, bestCount uint32
  bestQuality := float64(count)

  for i := uint32(0); i < 64; i++ {
    currMask := And(GetBits256FromTMX(traits, i), mask)
    currCount := currMask.Count()
    currQuality := math.Abs((float64(count) / 2) - float64(currCount))
    if (currQuality < bestQuality) {
      bestQuality = currQuality
      bestIndex = i
      bestCount = currCount
      bestMask = currMask
    }

    if (bestQuality < 0.75) {
      return bestMask, bestIndex, bestCount, bestQuality
    }
  }
  return bestMask, bestIndex, bestCount, bestQuality
}








func createLeaf( index uint32, mask Bits256, dt *DecTree, tags []string ) {
  tag := tags[mask.OnlyIndex()]
  rule := "F_" + tag
  dt.nodes[index] = DTNode{ rule, mask }
}

func createDeadlock( index, count uint32, mask Bits256, dt *DecTree,
                     tags []string ) {
  dlTags := []string{}
  for _, cid := range mask.AllIndices() {
    dlTags = append(dlTags, tags[cid])
  }
  //rule := fmt.Sprintf("DEADLOCK(%s)", strings.Join(dlTags, ";"))
  //dt.nodes[index] = DTNode{ rule, mask }
  //  fmt.Println("------------------->", rule)
  solveDeadlock( index, dlTags, tags, dt)
}

func createStandard( index, traitId uint32, mask Bits256, dt *DecTree ) {
  dt.nodes[index] = DTNode{ fmt.Sprintf("S_%02d", traitId), mask }
}

func createTiebreaker( index, tbId uint32, mask Bits256, dt *DecTree ) {
  dt.nodes[index] = DTNode{ fmt.Sprintf("T_%02d", tbId), mask }
}




