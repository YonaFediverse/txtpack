package txtpack

import (
	"bufio"
	"bytes"
	"github.com/itsabgr/go-handy"
	"github.com/jmcvetta/randutil"
	"testing"
)

func randUtf8Bytes(length int) []byte {
	str, err := randutil.AlphaString(length)
	handy.Throw(err)
	return []byte(str)
}
func randPair() *Pair {
	return NewPair(randUtf8Bytes(9), randUtf8Bytes(9))
}
func randPack(pairs uint) Pack {
	pack := make(Pack, 0)
	for range handy.N(pairs) {
		pack = append(pack, randPair())
	}
	return pack
}
func randPacks(n uint) []Pack {
	packs := make([]Pack, 0)
	for range handy.N(n) {
		packs = append(packs, randPack(9))
	}
	return packs
}
func TestCodec(t *testing.T) {
	packs := randPacks(10)
	buf := new(bytes.Buffer)
	err := NewEncoder(buf).Encode(packs...)
	if err != nil {
		t.Fatal(err)
	}
	packs2, err := NewDecoder(bufio.NewReader(buf)).Decode()
	if err != nil {
		t.Fatal(err)
	}
	for i := range packs2 {
		for j := range packs2[i] {
			if false == packs[i][j].Equal(packs2[i][j]) {
				t.Fail()
			}
		}
	}
}
