package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"sync"

	"github.com/new-world-tools/go-oodle"
)

var filecount uint32

type Arc struct {
	mu    sync.Mutex
	data  *os.File
	toc   *TOC
	files uint32
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

func (a *Arc) NewTOC(offset uint64) {
	a.toc = &TOC{
		offset: offset,
	}
}

func (a *Arc) WriteTOCEntries() {
	for i := 0; i < int(a.files); i++ {
		l := ReadUint32(a.data, 0)
		n := ReadString(a.data, l)
		o := ReadUint64(a.data, 0)
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

func (a *Arc) unpackEntry(i int, meta *os.File, verbose bool) {
	var next_offset uint64

	filecount++
	name := a.toc.Entries[i].name
	nameLen := a.toc.Entries[i].nameLen
	offset := a.toc.Entries[i].offset

	os.MkdirAll(path.Dir((name)), 0700)
	file, _ := os.Create(name)
	if i != len(a.toc.Entries)-1 {
		next_offset = a.toc.Entries[i+1].offset
	} else {
		next_offset = uint64(a.toc.offset)
	}
	a.mu.Lock()
	binary.Write(meta, binary.LittleEndian, nameLen)
	meta.WriteString(name)
	magic_comp := ReadUint32(a.data, int64(offset)) //"ozip"

	if magic_comp != 0x70697A6F {
		binary.Write(meta, binary.LittleEndian, int8(0))
		a.mu.Unlock()
		data_buf := make([]byte, (next_offset - offset))
		a.data.ReadAt(data_buf, int64(offset))
		file.Write(data_buf)
	} else {

		binary.Write(meta, binary.LittleEndian, int8(1))
		unk := ReadUint32(a.data, int64(offset)+4) //unk
		binary.Write(meta, binary.LittleEndian, unk)
		_ = ReadUint64(a.data, int64(offset)+8)  //dec_size
		_ = ReadUint32(a.data, int64(offset)+16) //"ozip"
		dec_size := ReadUint64(a.data, int64(offset)+20)
		c_size := ReadUint32(a.data, int64(offset)+28)
		a.mu.Unlock()
		data_buf := make([]byte, c_size)
		a.data.ReadAt(data_buf, int64(offset)+32)
		data, err := oodle.Decompress(data_buf, int64(dec_size))
		if err != nil {
			fmt.Println(err)
		}
		file.Write(data)
	}
	file.Close()
	if verbose {
		done := int(float64(filecount) / float64(a.files) * count)
		percent := math.Round(float64(filecount) / float64(a.files) * 100)
		bar := "[" + string(repeat('#', done)) + string(repeat(' ', count-done)) + "]"
		fmt.Printf("\r%s %v%%", bar, percent)
	}
}

func (a *Arc) repackEntry(meta *os.File, v int, verbose bool) {
	var (
		c_flag   = make([]byte, 1)
		name_buf bytes.Buffer
		offset   int64
	)
	a.mu.Lock()
	filecount++
	name_len := ReadUint32(meta, 0)
	io.CopyN(&name_buf, meta, int64(name_len))
	meta.Read(c_flag)
	file, _ := os.Open(name_buf.String())
	info, _ := file.Stat()
	file_size := info.Size()
	if c_flag[0] == 1 {
		var temp bytes.Buffer
		temp.WriteString("ozip")
		binary.Write(&temp, binary.LittleEndian, ReadUint32(meta, 0))
		binary.Write(&temp, binary.LittleEndian, file_size)
		temp.WriteString("ozip")
		binary.Write(&temp, binary.LittleEndian, file_size)
		a.mu.Unlock()
		f_buf := make([]byte, file_size)
		file.Read(f_buf)
		c_data, err := oodle.Compress(f_buf, oodle.CompressorKraken, v)
		if err != nil {
			fmt.Println(err)
		}
		a.mu.Lock()
		offset, _ = a.data.Seek(0, 1)
		a.WriteTOCEntry(name_len, name_buf.String(), uint64(offset))
		binary.Write(&temp, binary.LittleEndian, uint32(len(c_data)))
		a.data.Write(temp.Bytes())
		a.data.Write(c_data)
		a.mu.Unlock()
		f_buf = nil
		c_data = nil
	} else {
		offset, _ = a.data.Seek(0, 1)
		a.WriteTOCEntry(name_len, name_buf.String(), uint64(offset))
		file.WriteTo(a.data)
		a.mu.Unlock()
	}
	if verbose {
		done := int(float64(filecount) / float64(a.files) * count)
		percent := math.Round(float64(filecount) / float64(a.files) * 100)
		bar := "[" + string(repeat('#', done)) + string(repeat(' ', count-done)) + "]"
		fmt.Printf("\r%s %v%%", bar, percent)
	}
}

func ReadUint32(r io.Reader, pos int64) uint32 {
	buf := make([]byte, 4)
	f, ok := r.(*os.File)
	if ok && pos != 0 {
		f.ReadAt(buf, pos)
	}
	r.Read(buf)
	return binary.LittleEndian.Uint32(buf)
}

func ReadUint64(r io.Reader, pos int64) uint64 {
	buf := make([]byte, 8)
	f, ok := r.(*os.File)
	if ok && pos != 0 {
		f.ReadAt(buf, pos)
	}
	r.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

func ReadString(r io.Reader, len uint32) string {
	var out string
	b := make([]byte, 1)
	for i := 0; i < int(len); i++ {
		r.Read(b)
		out += string(b)
	}
	return out
}

func repeat(char rune, count int) []rune {
	res := make([]rune, count)
	for i := range res {
		res[i] = char
	}
	return res
}
