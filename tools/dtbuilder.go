package main

import (
  "os"
  "../pkg/flagger"
)

func main() {
  traitsPath := "../data/kb_traits.txt"
  dtPath := "../data/dt_nodes.txt"

  traitsFile, err := os.Open(traitsPath)
  if err != nil {
    panic(err)
  }
  defer traitsFile.Close()

  tm := flagger.LoadTraitMap(traitsFile)
  tmx := flagger.GetTraitMatrix(tm)
  dt := flagger.BuildDecTree(uint32(len(tm)), 1, &tmx, tm, 0.1)

  dtFile, err := os.Create(dtPath)
  if err != nil {
    panic(err)
  }
  defer dtFile.Close()

  dt.Write(dtFile)
}
