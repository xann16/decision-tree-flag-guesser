package flagger

import (
  "fmt"
)

type TraitMatrix [256][64]bool

func GetTraitMatrix( tm TraitMap ) ( tmx TraitMatrix ) {
  for _, n := range tm {
    if n.Id > 255 {
      panic(fmt.Errorf("Country Id exceeds 255 (%d)", n.Id))
    }
    for i := uint32(0); i < 64; i++ {
      if ( n.Traits & ( 1 << i ) ) != 0 {
        tmx[n.Id][i] = true
      }
    }
  }
  return
}





