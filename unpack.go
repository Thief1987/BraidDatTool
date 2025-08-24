package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func Unpack(a *os.File, threads int, verbose bool) {
	var (
		wg sync.WaitGroup
	)
	s := time.Now()
	meta, _ := os.Create("metadata.bin")
	arc := newArc(a)
	magic := ReadString(arc.data, 8)
	if magic != dat_magic {
		log.Fatalf("Invalid archive format, magic header not equal to %s", dat_magic)
	}
	arc.data.Seek(8, 0)
	TOC_offset := ReadUint64(arc.data, 0)
	arc.data.Seek(int64(TOC_offset), 0)
	arc.files = ReadUint32(arc.data, 0)
	binary.Write(meta, binary.LittleEndian, arc.files)
	arc.NewTOC(TOC_offset)
	arc.WriteTOCEntries()
	if threads == 1 {
		for i := 0; i < int(arc.files); i++ {
			arc.unpackEntry(i, meta, verbose)
		}
	} else {
		sem := make(chan struct{}, threads)
		wg.Add(int(arc.files))
		for i := 0; i < int(arc.files); i++ {
			sem <- struct{}{}
			go func(meta *os.File) {
				arc.unpackEntry(i, meta, verbose)
				wg.Done()
				<-sem
			}(meta)
		}
		wg.Wait()
	}
	f := time.Now()
	if verbose {
		fmt.Printf("\n\n%v files successfully unpacked in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
	}
	meta.Close()
}
