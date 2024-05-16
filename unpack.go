package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/new-world-tools/go-oodle"
)

func unpack(arc *os.File) {
	var (
		TOC_buf     bytes.Buffer
		next_offset uint64
		size        int
	)
	fmt.Println("Offset       Size                     Name   ")
	s := time.Now()
	filecount := 0
	meta, _ := os.Create("metadata.bin")
	arc.Seek(8, 0)
	info, _ := arc.Stat()
	arc_size := info.Size()
	TOC_offset := ReadUint32(arc)
	arc.Seek(int64(TOC_offset), 0)
	files := ReadUint32(arc)
	binary.Write(meta, binary.LittleEndian, files)
	io.CopyN(&TOC_buf, arc, arc_size-int64(TOC_offset))
	TOC_reader := bytes.NewReader(TOC_buf.Bytes())
	for i := 0; i < int(files); i++ {
		filecount++
		var buf, name_buf bytes.Buffer
		name_len := ReadUint32(TOC_reader)
		binary.Write(meta, binary.LittleEndian, name_len)
		io.CopyN(&name_buf, TOC_reader, int64(name_len))
		meta.WriteString(name_buf.String())
		os.MkdirAll(path.Dir((name_buf.String())), 0700)
		file, _ := os.Create(name_buf.String())
		offset := ReadUint64(TOC_reader)
		savepos, _ := TOC_reader.Seek(0, 1)
		if i != int(files-1) {
			next_name_len := ReadUint32(TOC_reader)
			TOC_reader.Seek(int64(next_name_len), 1)
			next_offset = ReadUint64(TOC_reader)
		} else {
			next_offset = uint64(TOC_offset)
		}
		arc.Seek(int64(offset), 0)
		magic_comp := ReadUint32(arc) //"ozip"

		if magic_comp != 0x70697A6F {
			binary.Write(meta, binary.LittleEndian, int8(0))
			arc.Seek(int64(offset), 0)
			io.CopyN(file, arc, int64(next_offset-offset))
			size = int(next_offset - offset)
		} else {
			binary.Write(meta, binary.LittleEndian, int8(1))
			unk := ReadUint32(arc) //unk
			binary.Write(meta, binary.LittleEndian, unk)
			_ = ReadUint64(arc) //dec_size
			_ = ReadUint32(arc) //"ozip"
			dec_size := ReadUint64(arc)
			c_size := ReadUint32(arc)
			io.CopyN(&buf, arc, int64(c_size))
			data, err := oodle.Decompress(buf.Bytes(), int64(dec_size))
			if err != nil {
				fmt.Println(err)
			}
			size = len(data)
			file.Write(data)
		}
		file.Close()
		TOC_reader.Seek(savepos, 0)
		fmt.Printf("0x%X       %v        %s\n", offset, size, name_buf.String())

	}
	f := time.Now()
	fmt.Printf("%v files succesfully unpacked in %.2f sec", filecount, f.Sub(s).Abs().Seconds())
	meta.Close()

}
