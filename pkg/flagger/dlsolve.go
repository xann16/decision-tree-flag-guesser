package flagger

import (
  "fmt"
)

func hasTag(tags []string, tag string) bool {
  for _, t := range tags {
    if t == tag {
      return true
    }
  }
  return false
}

func findIndex(tags []string, tag string) uint32 {
    for i, t := range tags {
    if t == tag {
      return uint32(i)
    }
  }
  panic(fmt.Errorf("tag not found"))
}



func solveDeadlock( index uint32, currTags, tags []string, dt *DecTree ) {
  switch len(currTags) {
  case 4:
    solveLUNLRSRU(index, currTags, tags, dt)
  case 3:
    if hasTag(currTags, "id") {
      solveIDMCPL(index, currTags, tags, dt)
    } else if hasTag(currTags, "ch") {
      solveCHENGE(index, currTags, tags, dt)
    } else if hasTag(currTags, "fo") {
      solveFOISNO(index, currTags, tags, dt)
    }
  case 2:
    if hasTag(currTags, "bo") {
      solvePair( index, 1, "bo", "lt", tags, dt)
    } else if hasTag(currTags, "ps") {
      solvePair( index, 1, "sd", "ps", tags, dt)
    } else if hasTag(currTags, "bg") {
      solvePair( index, 1, "hu", "bg", tags, dt)
    } else if hasTag(currTags, "gn") {
      solvePair( index, 7, "gn", "ml", tags, dt)
    } else if hasTag(currTags, "ci") {
      solvePair( index, 7, "ci", "ie", tags, dt)
    } else if hasTag(currTags, "gn") {
      solvePair( index, 0, "th", "cr", tags, dt)
    } else if hasTag(currTags, "td") {
      solvePair( index, 8, "td", "ro", tags, dt)
    } else if hasTag(currTags, "cu") {
      solvePair( index, 0, "cu", "pr", tags, dt)
    }
  }

}

/*
[][] TIEBREAKERS:

[ 0] blue stripe in the middle
[ 1] red stripe on the top
[ 2] has colour: light blue
[ 3] has aspect ratio 4:5 (squarish)
[ 4] cross is red
[ 5] has more than one cross
[ 6] has blue background
[ 7] green stripe is on the right
[ 8] blue is slightly darker

*/


func solveLUNLRSRU( index uint32, currTags, tags []string, dt *DecTree ) {
  lu := GetSingle256(findIndex(tags, "lu"))
  nl := GetSingle256(findIndex(tags, "nl"))
  rs := GetSingle256(findIndex(tags, "rs"))
  ru := GetSingle256(findIndex(tags, "ru"))

  createTiebreaker(index, 0, Or(Or(Or(lu, nl), rs), ru), dt)

  index <<= 1

  createTiebreaker(index, 1, Or(rs, ru), dt)
  createTiebreaker(index + 1, 2, Or(nl, lu), dt)

  index <<= 1

  createLeaf(index + 0, rs, dt, tags)
  createLeaf(index + 1, ru, dt, tags)
  createLeaf(index + 2, lu, dt, tags)
  createLeaf(index + 3, nl, dt, tags)
}

func solveIDMCPL( index uint32, currTags, tags []string, dt *DecTree ) {
  id := GetSingle256(findIndex(tags, "id"))
  mc := GetSingle256(findIndex(tags, "mc"))
  pl := GetSingle256(findIndex(tags, "pl"))

  createTiebreaker(index, 1, Or(Or(id, mc), pl), dt)

  index <<= 1

  createTiebreaker(index, 3, Or(id, mc), dt)
  createLeaf(index + 1, pl, dt, tags)

  index <<= 1

  createLeaf(index + 0, mc, dt, tags)
  createLeaf(index + 1, id, dt, tags)
}

func solveCHENGE( index uint32, currTags, tags []string, dt *DecTree ) {
  ch := GetSingle256(findIndex(tags, "ch"))
  en := GetSingle256(findIndex(tags, "en"))
  ge := GetSingle256(findIndex(tags, "ge"))

  createTiebreaker(index, 4, Or(Or(ch, en), ge), dt)

  index <<= 1

  createTiebreaker(index, 5, Or(en, ge), dt)
  createLeaf(index + 1, ch, dt, tags)

  index <<= 1

  createLeaf(index + 0, ge, dt, tags)
  createLeaf(index + 1, en, dt, tags)
}

func solveFOISNO( index uint32, currTags, tags []string, dt *DecTree ) {
  fo := GetSingle256(findIndex(tags, "fo"))
  is := GetSingle256(findIndex(tags, "is"))
  no := GetSingle256(findIndex(tags, "no"))

  createTiebreaker(index, 4, Or(Or(fo, is), no), dt)

  index <<= 1

  createTiebreaker(index, 6, Or(is, fo), dt)
  createLeaf(index + 1, no, dt, tags)

  index <<= 1

  createLeaf(index + 0, is, dt, tags)
  createLeaf(index + 1, fo, dt, tags)
}

func solvePair( index, tbId uint32, yesTag, noTag string,
                tags []string, dt *DecTree ) {
  yt := GetSingle256(findIndex(tags, yesTag))
  nt := GetSingle256(findIndex(tags, noTag))

  createTiebreaker(index, tbId, Or(yt, nt), dt)

  index <<= 1

  createLeaf(index + 0, yt, dt, tags)
  createLeaf(index + 1, nt, dt, tags)
}
