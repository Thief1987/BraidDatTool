package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const dat_magic = "BRAID-BF"

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

func main() {

	//os.Mkdir(f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))], 0700)
	//os.Chdir(f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))])

	//unpack(f)

	args := os.Args
	if len(args) == 1 {
		fmt.Printf("Usage : \n      unpack: -u archive_name \n      repack: -r archive name compression_level\n")
		fmt.Printf("When repacking, optionally you can specify compression level, legit values are from -4(fastest) to 9(slowest).\n")
		fmt.Printf("Default value is 6(devs used it), but it's pretty slow, very slow I would say, so I decided to add this option at least for the testing purposes.")
		log.Fatal()
	}
	
	if args[1] == "-u" {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		unpack(f)
	} else if args[1] == "-r" {
		if len(args) > 3 {
			value, _ := strconv.Atoi(args[3])
			if value >= -4 && value <= 9 {
				fmt.Printf("Compression level value is set to %v\n", value)
				repack(value)
			} else {
				fmt.Printf("Invalid compression level value. Value will be set to 6")
				value = 6
				repack(value)
			}
		} else {
			fmt.Printf("Compression level value is not specified. Value will be set to 6")
			value := 6
			repack(value)
		}

	} else {
		fmt.Printf("Usage : \n      unpack: -u archive_name \n      repack: -r archive name compression_level\n")
		fmt.Printf("When repacking optionally you can specify compression level, legit values are from -4(fastest) to 9(slowest).\n")
		fmt.Printf("Default value is 6(devs used it), but it's pretty slow, very slow I would say, so I decided to add this option at least for testing purposes.")
		log.Fatal()
	}
}
