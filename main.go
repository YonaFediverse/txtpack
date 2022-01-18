package txtpack

import (
	"bufio"
	"bytes"
	"io"
)

func SplitLine(line []byte) (pair *Pair) {
	pair = new(Pair)
	parts := bytes.SplitN(line, []byte(":"), 2)
	pair.key = bytes.TrimSpace(parts[0])
	if len(parts) == 2 {
		pair.val = bytes.TrimSpace(parts[1])
	}
	if len(pair.key) == 0 {
		pair.key = nil
	}
	if len(pair.val) == 0 {
		pair.val = nil
	}
	return pair
}

type Pair struct {
	key []byte
	val []byte
}

func (p *Pair) Equal(other *Pair) bool {
	return bytes.Equal(p.val, other.Value()) && bytes.Equal(p.key, other.Key())
}
func (p *Pair) Join() []byte {
	b := new(bytes.Buffer)
	b.Write(p.key)
	b.WriteByte(':')
	b.Write(p.val)
	return b.Bytes()
}
func NewPair(key, val []byte) *Pair {
	return &Pair{
		key: bytes.TrimSpace(key),
		val: bytes.TrimSpace(val),
	}
}
func (p *Pair) IsNil() bool {
	return len(p.key) == 0 && len(p.val) == 0
}
func (p *Pair) IsComment() bool {
	return len(p.key) == 0 && len(p.val) > 0
}
func (p *Pair) Value() []byte {
	return p.val
}

func (p *Pair) Key() []byte {
	return p.key
}

type Pack []*Pair

func (p *Pair) String() string {
	return string(p.Join())
}
func (p Pack) String() string {
	return string(p.Encode())
}
func (p Pack) Encode() []byte {
	b := new(bytes.Buffer)
	for _, pair := range p {
		b.Write(pair.Join())
		b.WriteByte('\n')
	}
	return b.Bytes()
}

type Encoder struct {
	dst io.Writer
}

func NewEncoder(dst io.Writer) *Encoder {
	return &Encoder{dst: dst}
}
func NewPack(pairs ...*Pair) Pack {
	return pairs
}
func (e *Encoder) Writeln(line []byte) error {
	_, err := e.dst.Write(bytes.TrimSpace(line))
	if err != nil {
		return err
	}
	_, err = e.dst.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}
func (e *Encoder) WritePair(p *Pair) error {
	return e.Writeln(p.Join())
}
func (e *Encoder) WritePack(pack Pack) error {
	return e.Writeln(pack.Encode())
}

func (e *Encoder) Encode(packs ...Pack) error {
	for _, pack := range packs {
		err := e.WritePack(pack)
		if err != nil {
			return err
		}
	}
	return nil
}

type Decoder struct {
	src *bufio.Reader
}

func NewDecoder(src *bufio.Reader) *Decoder {
	return &Decoder{src: src}
}
func (d *Decoder) Decode() ([]Pack, error) {
	packs := make([]Pack, 0)
	for {
		pack, err := d.NextPack()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		packs = append(packs, pack)
	}
	return packs, nil
}

func (d *Decoder) NextPack() (Pack, error) {
	pack := make(Pack, 0)
	for {
		pair, err := d.NextPair()
		if err != nil {
			return nil, err
		}
		if pair.IsNil() {
			break
		}
		pack = append(pack, pair)
	}
	return pack, nil
}
func (d *Decoder) NextPair() (*Pair, error) {
	line, err := d.NextLine()
	if err != nil {
		return nil, err
	}
	return SplitLine(line), nil
}

func (d *Decoder) NextLine() ([]byte, error) {
	b := new(bytes.Buffer)
	for {
		line, isPrefix, err := d.src.ReadLine()
		if err != nil {
			return nil, err
		}
		b.Write(line)
		if false == isPrefix {
			break
		}
	}
	return b.Bytes(), nil
}
