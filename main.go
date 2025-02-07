package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	dat_magic     = "BRAID-BF"
	threadsUnpack = 16
	threadsPack   = 50
)

func die_with_usage_message() {
	fmt.Printf("Usage: BraidDatTool [-u archive_name | -r archive_name [compression_level]]\n")
	fmt.Printf("\t -u archive_name                            Unpack the archive. \n")
	fmt.Printf("\t -r archive_name [compression_level]        Repack the archive. \n")
	fmt.Printf("When repacking, optionally you can specify the compression level. Valid values are from -4 (fastest) to 9 (slowest).\n")
	fmt.Printf("Default value is 6 (devs used it), but it's pretty slow, very slow I would say, so I decided to add this compression level option at least for testing purposes.\n")
	fmt.Printf("Looking for the archive to run BraidDatTool on? Maybe its \"C:\\Program Files (x86)\\Steam\\steamapps\\common\\Braid Anniversary Edition\\data\\data.dat\" or in a similar location.\n")
	log.Fatal()
}

func dllDump() {
	dll, _ := os.Create("oo2core_9_win64.dll")
	dll.Write(oo2core_9_win64_dll)
	dll.Close()
}

func main() {
	_, err := os.Stat("oo2core_9_win64.dll")
	if err != nil {
		if os.IsNotExist(err) {
			dllDump()
		}
	}

	args := os.Args
	if len(args) == 1 {
		die_with_usage_message()
	}

	if args[1] == "-u" {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		Unpack(f)
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
		die_with_usage_message()
	}
}
