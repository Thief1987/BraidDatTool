package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func Repack(v int, threadsPack uint16) {
	var (
		wg sync.WaitGroup
	)
	fmt.Println("Offset       Size                     Name   ")
	s := time.Now()
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
	if threadsPack == 1 {
		for i := 0; i < int(arc.files); i++ {
			arc.repackEntry(meta, v)
		}
	} else {
		wg.Add(int(arc.files))
		sem := make(chan struct{}, threadsPack)
		for i := 0; i < int(arc.files); i++ {
			sem <- struct{}{}
			go func(meta *os.File) {
				arc.repackEntry(meta, v)
				wg.Done()
				<-sem
			}(meta)
		}
		wg.Wait()
	}
	offset, _ := new_arc.Seek(0, 1)
	arc.toc.offset = uint64(offset)
	arc.WriteTOCBinary()
	arc.data.Seek(8, 0)
	binary.Write(arc.data, binary.LittleEndian, arc.toc.offset)
	arc.data.Close()
	f := time.Now()
	fmt.Printf("%v files successfully packed in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
}
