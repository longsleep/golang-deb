package main

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func rnd(v int32, r int32) int32 {
	if r <= 0 {
		return v
	}
	v += r - 1
	c := v % r
	if c < 0 {
		c += r
	}
	v -= c
	return v
}

func readwithpad(r io.Reader, sz int32) ([]byte, error) {
	full := rnd(sz, 4)
	data := make([]byte, full)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	data = data[:sz]
	return data, nil
}

func readnote(filename, name string, type_ int32) ([]byte, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return nil, err
	}
	for _, sect := range f.Sections {
		if sect.Type != elf.SHT_NOTE {
			continue
		}
		r := sect.Open()
		for {
			var namesize, descsize, nt_type int32
			err = binary.Read(r, f.ByteOrder, &namesize)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, fmt.Errorf("read namesize failed", err)
			}
			err = binary.Read(r, f.ByteOrder, &descsize)
			if err != nil {
				return nil, fmt.Errorf("read descsize failed", err)
			}
			err = binary.Read(r, f.ByteOrder, &nt_type)
			if err != nil {
				return nil, fmt.Errorf("read type failed", err)
			}
			nt_name, err := readwithpad(r, namesize)
			if err != nil {
				return nil, fmt.Errorf("read name failed", err)
			}
			desc, err := readwithpad(r, descsize)
			if err != nil {
				return nil, fmt.Errorf("read desc failed", err)
			}
			if name == string(nt_name) && type_ == nt_type {
				return desc, nil
			}
		}
	}
	return nil, nil
}

func main() {
	desc, err := readnote(os.Args[1], "Go\x00\x00", 2)
	if err != nil {
		log.Fatalf("readnote failed: %v", err)
	}
	fmt.Printf("%x\n", desc)
}
