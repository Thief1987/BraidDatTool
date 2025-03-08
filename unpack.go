package main

import (
	"encoding/binary"
	"os"
	"sync"
)

func Unpack(a *os.File, threadsUnpack uint16) {
	var (
		wg, wg1 sync.WaitGroup
	)
	//fmt.Println("Offset       Size                     Name   ")
	//s := time.Now()
	meta, _ := os.Create("metadata.bin")
	arc := newArc(a)
	arc.data.Seek(8, 0)
	TOC_offset := ReadUint64(arc.data)
	arc.data.Seek(int64(TOC_offset), 0)
	arc.files = ReadUint32(arc.data)
	binary.Write(meta, binary.LittleEndian, arc.files)
	arc.NewTOC(TOC_offset)
	arc.WriteTOCEntries()
	for i := 0; i < int(arc.files/uint32(threadsUnpack)); i++ {
		wg.Add(1)
		for j := 0; j < int(threadsUnpack); j++ {
			wg1.Add(1)
			go func(meta *os.File) {
				arc.unpackEntry((i*int(threadsUnpack) + j), meta)
				wg1.Done()
			}(meta)
		}
		wg1.Wait()
		wg.Done()
	}
	wg.Wait()
	remain := int(arc.files % uint32(threadsUnpack))
	for i := 0; i < remain; i++ {
		arc.unpackEntry(int(arc.files)-remain+i, meta)
	}
	//f := time.Now()
	//fmt.Printf("%v files successfully unpacked in %.2f sec\n", filecount, f.Sub(s).Abs().Seconds())
	meta.Close()
}
