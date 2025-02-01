package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/new-world-tools/go-oodle"
)

func repack(v int) {
	var TOC_buf bytes.Buffer
	fmt.Println("Offset       Size                     Name   ")
	s := time.Now()
	filecount := 0
	new_arc, _ := os.Create("braid.dat_new")
	meta, err := os.Open("metadata.bin")
	if err != nil {
		new_arc.Close()
		os.Remove("braid.dat_new")
		log.Fatal("metadata.bin doesn't exist; try to unpack the original archive first")
	}
	new_arc.WriteString(dat_magic)
	binary.Write(new_arc, binary.LittleEndian, uint64(0))
	files := ReadUint32(meta)
	binary.Write(&TOC_buf, binary.LittleEndian, files)

	for i := 0; i < int(files); i++ {
		var (
			c_flag   = make([]byte, 1)
			name_buf bytes.Buffer
			size     int
		)
		filecount++
		name_len := ReadUint32(meta)
		binary.Write(&TOC_buf, binary.LittleEndian, name_len)
		io.CopyN(&name_buf, meta, int64(name_len))
		TOC_buf.WriteString(name_buf.String())
		offset, _ := new_arc.Seek(0, 1)
		binary.Write(&TOC_buf, binary.LittleEndian, offset)
		meta.Read(c_flag)
		file, _ := os.Open(name_buf.String())
		info, _ := file.Stat()
		file_size := info.Size()
		if c_flag[0] == 1 {
			var f_buf bytes.Buffer
			new_arc.WriteString("ozip")
			binary.Write(new_arc, binary.LittleEndian, ReadUint32(meta))
			binary.Write(new_arc, binary.LittleEndian, file_size)
			new_arc.WriteString("ozip")
			binary.Write(new_arc, binary.LittleEndian, file_size)
			io.Copy(&f_buf, file)
			c_data, _ := oodle.Compress(f_buf.Bytes(), oodle.CompressorKraken, v)
			size = len(c_data)
			binary.Write(new_arc, binary.LittleEndian, uint32(len(c_data)))
			new_arc.Write(c_data)
		} else {
			size = int(file_size)
			io.Copy(new_arc, file)
		}
		fmt.Printf("0x%X       %v        %s\n", offset, size, name_buf.String())
	}
	TOC_offset, _ := new_arc.Seek(0, 1)
	new_arc.Write(TOC_buf.Bytes())
	new_arc.Seek(8, 0)
	binary.Write(new_arc, binary.LittleEndian, TOC_offset)
	new_arc.Close()
	f := time.Now()
	fmt.Printf("%v files successfully packed in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
}
