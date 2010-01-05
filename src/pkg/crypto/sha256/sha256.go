// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements the SHA256 hash algorithm as defined in FIPS 180-2.
package sha256

import (
	"hash"
	"os"
)

// The size of a SHA256 checksum in bytes.
const Size = 32

const (
	_Chunk = 64
	_Init0 = 0x6A09E667
	_Init1 = 0xBB67AE85
	_Init2 = 0x3C6EF372
	_Init3 = 0xA54FF53A
	_Init4 = 0x510E527F
	_Init5 = 0x9B05688C
	_Init6 = 0x1F83D9AB
	_Init7 = 0x5BE0CD19
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	h   [8]uint32
	x   [_Chunk]byte
	nx  int
	len uint64
}

func (d *digest) Reset() {
	d.h[0] = _Init0
	d.h[1] = _Init1
	d.h[2] = _Init2
	d.h[3] = _Init3
	d.h[4] = _Init4
	d.h[5] = _Init5
	d.h[6] = _Init6
	d.h[7] = _Init7
	d.nx = 0
	d.len = 0
}

// New returns a new hash.Hash computing the SHA256 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) Write(p []byte) (nn int, err os.Error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := len(p)
		if n > _Chunk-d.nx {
			n = _Chunk - d.nx
		}
		for i := 0; i < n; i++ {
			d.x[d.nx+i] = p[i]
		}
		d.nx += n
		if d.nx == _Chunk {
			_Block(d, &d.x)
			d.nx = 0
		}
		p = p[n:]
	}
	n := _Block(d, p)
	p = p[n:]
	if len(p) > 0 {
		for i, x := range p {
			d.x[i] = x
		}
		d.nx = len(p)
	}
	return
}

func (d *digest) Sum() []byte {
	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	len := d.len
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 {
		d.Write(tmp[0 : 56-len%64])
	} else {
		d.Write(tmp[0 : 64+56-len%64])
	}

	// Length in bits.
	len <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(len >> (56 - 8*i))
	}
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panicln("oops")
	}

	p := make([]byte, 32)
	j := 0
	for _, s := range d.h {
		p[j+0] = byte(s >> 24)
		p[j+1] = byte(s >> 16)
		p[j+2] = byte(s >> 8)
		p[j+3] = byte(s >> 0)
		j += 4
	}
	return p
}