package flagger

import (
	"fmt"
	"math/bits"
	"strconv"
)

const ONES64 = 0xFFFFFFFFFFFFFFFF

type Bits256 struct {
	d [4]uint64
}

func (bs *Bits256) SetOne(pos uint32) {
	if pos > 255 {
		panic(fmt.Errorf("Index out of bounds"))
	}
	bs.d[pos>>6] |= (1 << (pos & 0x3F))
}

func (bs *Bits256) SetZero(pos uint32) {
	if pos > 255 {
		panic(fmt.Errorf("Index out of bounds"))
	}
	bs.d[pos>>6] &^= (uint64(1) << (pos & 0x3F))
}

func (bs *Bits256) Set(pos uint32, val bool) {
	if val {
		bs.SetOne(pos)
	} else {
		bs.SetZero(pos)
	}
}

func (bs Bits256) Get(pos uint32) bool {
	return (bs.d[pos>>6] & (uint64(1) << (pos & 0x3F))) != 0
}

func (bs Bits256) Count() uint32 {
	return uint32(bits.OnesCount64(bs.d[0]) + bits.OnesCount64(bs.d[1]) +
		bits.OnesCount64(bs.d[2]) + bits.OnesCount64(bs.d[3]))
}

func (bs Bits256) OnlyIndex() uint32 {
	if bs.Count() != 1 {
		panic(fmt.Errorf("More than one flag set"))
	}
	for i := uint32(0); i < 256; i++ {
		if bs.Get(i) {
			return i
		}
	}
	panic(fmt.Errorf("Flag not found"))
}

func (bs Bits256) AllIndices() []uint32 {
	res := []uint32{}
	for i := uint32(0); i < 256; i++ {
		if bs.Get(i) {
			res = append(res, i)
		}
	}
	return res
}

func (bs Bits256) ToString() string {
	return fmt.Sprintf("%016x%016x%016x%016x", bs.d[3], bs.d[2], bs.d[1], bs.d[0])
}

func And(bs1, bs2 Bits256) (bs Bits256) {
	bs.d[0] = bs1.d[0] & bs2.d[0]
	bs.d[1] = bs1.d[1] & bs2.d[1]
	bs.d[2] = bs1.d[2] & bs2.d[2]
	bs.d[3] = bs1.d[3] & bs2.d[3]
	return
}

func Or(bs1, bs2 Bits256) (bs Bits256) {
	bs.d[0] = bs1.d[0] | bs2.d[0]
	bs.d[1] = bs1.d[1] | bs2.d[1]
	bs.d[2] = bs1.d[2] | bs2.d[2]
	bs.d[3] = bs1.d[3] | bs2.d[3]
	return
}

func Not(bsi Bits256) (bs Bits256) {
	bs.d[0] = ^bsi.d[0]
	bs.d[1] = ^bsi.d[1]
	bs.d[2] = ^bsi.d[2]
	bs.d[3] = ^bsi.d[3]
	return
}

func GetBits256(data ...uint64) (bs Bits256) {
	if len(data) > 4 {
		panic(fmt.Errorf("Too many params"))
	}
	for i, dp := range data {
		bs.d[i] = dp
	}
	return
}

func GetOnes256(count uint32) Bits256 {
	if count > 255 {
		panic(fmt.Errorf("Index out of bounds"))
	}
	res := GetBits256(ONES64, ONES64, ONES64, ONES64)
	for i := count; i < 256; i++ {
		res.SetZero(i)
	}
	return res
}

func GetSingle256(pos uint32) Bits256 {
	if pos > 255 {
		panic(fmt.Errorf("Index out of bounds"))
	}
	res := Bits256{}
	res.SetOne(pos)
	return res
}

func GetBits256FromString(str string) (bs Bits256) {
	if len(str) != 64 {
		panic(fmt.Errorf("Invalid length of bitset256 hex string (%d).", len(str)))
	}
	for i := 0; i < 4; i++ {
		start := i * 16
		end := (i + 1) * 16
		val, err := strconv.ParseUint(str[start:end], 16, 64)
		if err != nil {
			panic(err)
		}
		bs.d[3-i] = val
	}
	return
}

func GetBits256FromTMX(tmx *TraitMatrix, traitId uint32) (bs Bits256) {
	for i := uint32(0); i < 256; i++ {
		if tmx[i][traitId] {
			bs.SetOne(i)
		}
	}
	return
}
