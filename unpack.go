package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"
)

func Unpack(a *os.File, threadsUnpack uint16) {
	var (
		wg sync.WaitGroup
	)
	//fmt.Println("Offset       Size                     Name   ")
	s := time.Now()
	meta, _ := os.Create("metadata.bin")
	arc := newArc(a)
	arc.data.Seek(8, 0)
	TOC_offset := ReadUint64(arc.data)
	arc.data.Seek(int64(TOC_offset), 0)
	arc.files = ReadUint32(arc.data)
	binary.Write(meta, binary.LittleEndian, arc.files)
	arc.NewTOC(TOC_offset)
	arc.WriteTOCEntries()
	if threadsUnpack == 1 {
		for i := 0; i < int(arc.files); i++ {
			arc.unpackEntry(i, meta)
		}
	} else {
		sem := make(chan struct{}, threadsUnpack)
		wg.Add(int(arc.files))
		for i := 0; i < int(arc.files); i++ {
			sem <- struct{}{}
			go func(meta *os.File) {
				arc.unpackEntry(i, meta)
				wg.Done()
				<-sem
			}(meta)
		}
		wg.Wait()
	}
	f := time.Now()
	fmt.Printf("\r%v files successfully unpacked in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
	meta.Close()
}
