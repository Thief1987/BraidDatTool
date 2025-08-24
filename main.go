package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	_ "embed"
)

const (
	dat_magic = "BRAID-BF"
	count     = 50
)

//go:embed README.md
var readme string

func die_with_usage_message() {
	fmt.Printf(strings.Join(strings.Split(readme, "\n")[2:], "\n"))
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
	threads := runtime.NumCPU()

	switch args[1] {
	case "-u":
		f, _ := os.Open(args[2])
		defer f.Close()
		Unpack(f, threads, true)
	case "-r":
		value, _ := strconv.Atoi(args[3])
		arcName := args[2] + "_new"
		if len(args) > 3 {
			if value >= -4 && value <= 9 {
				fmt.Printf("Compression level value is set to %v\n\n", value)
				Repack(value, threads, arcName, true)
			} else {
				fmt.Printf("Invalid compression level value. Value will be set to 6.\n\n")
				value = 6
				Repack(value, threads, arcName, true)
			}
		} else {
			fmt.Printf("Compression level value is not specified. Value will be set to 6.\n\n")
			value := 6
			Repack(value, threads, arcName, true)
		}

	default:
		die_with_usage_message()
	}
}
