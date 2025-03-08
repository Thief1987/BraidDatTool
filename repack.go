package main

import (
	"encoding/binary"
	"log"
	"os"
	"sync"
)

func Repack(v int, threadsPack uint16) {
	var (
		wg, wg1 sync.WaitGroup
	)
	//fmt.Println("Offset       Size                     Name   ")
	//s := time.Now()
	new_arc, _ := os.Create("braid.dat_new")
	meta, err := os.Open("metadata.bin")
	if err != nil {
		new_arc.Close()
		os.Remove("braid.dat_new")
		log.Fatal("metadata.bin doesn't exist; try to unpack the original archive first")
	}
	arc := newArc(new_arc)
	arc.NewTOC(0)
	arc.data.WriteString(dat_magic)
	binary.Write(arc.data, binary.LittleEndian, uint64(0))
	arc.files = ReadUint32(meta)
	for i := 0; i < int(arc.files/uint32(threadsPack)); i++ {
		wg.Add(1)
		for j := 0; j < int(threadsPack); j++ {
			wg1.Add(1)
			go func(meta *os.File) {
				arc.repackEntry(meta, v)
				wg1.Done()
			}(meta)
		}
		wg1.Wait()
		wg.Done()
	}
	wg.Wait()
	remain := int(arc.files % uint32(threadsPack))
	for i := 0; i < remain; i++ {
		arc.repackEntry(meta, v)
	}
	offset, _ := new_arc.Seek(0, 1)
	arc.toc.offset = uint64(offset)
	arc.WriteTOCBinary()
	arc.data.Seek(8, 0)
	binary.Write(arc.data, binary.LittleEndian, arc.toc.offset)
	arc.data.Close()
	//f := time.Now()
	//fmt.Printf("%v files successfully packed in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
}
