package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/new-world-tools/go-oodle"
)

var filecount uint32

type Arc struct {
	mu    sync.Mutex
	files uint32
	data  *os.File
	toc   *TOC
}

type TOCEntry struct {
	nameLen uint32
	name    string
	offset  uint64
}

type TOC struct {
	offset  uint64
	Entries []*TOCEntry
}

func newArc(d *os.File) Arc {
	return Arc{
		data: d,
	}
}

func (a *Arc) ReadString(len uint32) string {
	var out string
	for i := 0; i < int(len); i++ {
		b := make([]byte, 1)
		a.data.Read(b)
		out += string(b)
	}
	return out
}

func (a *Arc) NewTOC(offset uint64) {
	a.toc = &TOC{
		offset: offset,
	}
}

func (a *Arc) WriteTOCEntries() {
	for i := 0; i < int(a.files); i++ {
		l := ReadUint32(a.data)
		n := a.ReadString(l)
		o := ReadUint64(a.data)
		entry := TOCEntry{
			nameLen: l,
			name:    n,
			offset:  o,
		}
		a.toc.Entries = append(a.toc.Entries, &entry)
	}
}

func (a *Arc) WriteTOCEntry(l uint32, n string, o uint64) {
	entry := TOCEntry{
		nameLen: l,
		name:    n,
		offset:  o,
	}
	a.toc.Entries = append(a.toc.Entries, &entry)
}

func (a *Arc) WriteTOCBinary() {
	binary.Write(a.data, binary.LittleEndian, a.files)
	for i := 0; i < int(a.files); i++ {
		binary.Write(a.data, binary.LittleEndian, a.toc.Entries[i].nameLen)
		a.data.WriteString(a.toc.Entries[i].name)
		binary.Write(a.data, binary.LittleEndian, a.toc.Entries[i].offset)
	}
}

func (a *Arc) unpackEntry(i int, meta *os.File) {
	var (
		next_offset uint64
		size        int
	)
	a.mu.Lock()
	filecount++
	name := a.toc.Entries[i].name
	nameLen := a.toc.Entries[i].nameLen
	offset := a.toc.Entries[i].offset
	var buf bytes.Buffer
	binary.Write(meta, binary.LittleEndian, nameLen)
	meta.WriteString(name)
	os.MkdirAll(path.Dir((name)), 0700)
	file, _ := os.Create(name)
	if i != len(a.toc.Entries)-1 {
		next_offset = a.toc.Entries[i+1].offset
	} else {
		next_offset = uint64(a.toc.offset)
	}
	a.data.Seek(int64(offset), 0)
	magic_comp := ReadUint32(a.data) //"ozip"

	if magic_comp != 0x70697A6F {
		binary.Write(meta, binary.LittleEndian, int8(0))
		a.data.Seek(int64(offset), 0)
		io.CopyN(file, a.data, int64(next_offset-offset))
		size = int(next_offset - offset)
		a.mu.Unlock()
	} else {
		binary.Write(meta, binary.LittleEndian, int8(1))
		unk := ReadUint32(a.data) //unk
		binary.Write(meta, binary.LittleEndian, unk)
		_ = ReadUint64(a.data) //dec_size
		_ = ReadUint32(a.data) //"ozip"
		dec_size := ReadUint64(a.data)
		c_size := ReadUint32(a.data)
		io.CopyN(&buf, a.data, int64(c_size))
		a.mu.Unlock()
		data, err := oodle.Decompress(buf.Bytes(), int64(dec_size))
		if err != nil {
			fmt.Println(err)
		}
		size = len(data)
		file.Write(data)
	}
	file.Close()
	fmt.Printf("0x%X       %v        %s\n", offset, size, name)
}

func (a *Arc) repackEntry(meta *os.File, v int) {
	var (
		c_flag   = make([]byte, 1)
		name_buf bytes.Buffer
		size     int
		offset   int64
	)
	a.mu.Lock()
	filecount++
	name_len := ReadUint32(meta)
	io.CopyN(&name_buf, meta, int64(name_len))
	meta.Read(c_flag)
	file, _ := os.Open(name_buf.String())
	info, _ := file.Stat()
	file_size := info.Size()
	if c_flag[0] == 1 {
		var f_buf, temp bytes.Buffer
		temp.WriteString("ozip")
		binary.Write(&temp, binary.LittleEndian, ReadUint32(meta))
		binary.Write(&temp, binary.LittleEndian, file_size)
		temp.WriteString("ozip")
		binary.Write(&temp, binary.LittleEndian, file_size)
		a.mu.Unlock()
		io.Copy(&f_buf, file)
		c_data, err := oodle.Compress(f_buf.Bytes(), oodle.CompressorKraken, v)
		if err != nil {
			fmt.Println(err)
		}
		size = len(c_data)
		a.mu.Lock()
		offset, _ = a.data.Seek(0, 1)
		a.WriteTOCEntry(name_len, name_buf.String(), uint64(offset))
		binary.Write(&temp, binary.LittleEndian, uint32(len(c_data)))
		a.data.Write(temp.Bytes())
		a.data.Write(c_data)
		a.mu.Unlock()
	} else {
		size = int(file_size)
		offset, _ = a.data.Seek(0, 1)
		a.WriteTOCEntry(name_len, name_buf.String(), uint64(offset))
		io.Copy(a.data, file)
		a.mu.Unlock()
	}
	fmt.Printf("0x%X       %v        %s\n", offset, size, name_buf.String())
}

func ReadUint32(r io.Reader) uint32 {
	var buf bytes.Buffer
	io.CopyN(&buf, r, 4)
	return binary.LittleEndian.Uint32(buf.Bytes())
}

func ReadUint64(r io.Reader) uint64 {
	var buf bytes.Buffer
	io.CopyN(&buf, r, 8)
	return binary.LittleEndian.Uint64(buf.Bytes())
}
