package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/anurag925/bloomcheck/src/filter"
)

func main() {
	// Define a flag for the -build option
	buildFlag := flag.Bool("build", false, "Build flag")

	// Parse the command-line flags
	flag.Parse()

	// Check if the -build flag is provided
	if *buildFlag {
		if err := buildBloomFilterHash(); err != nil {
			fmt.Println("unable to build the filter error: ", err)
			return
		}
		return
	}

	// Get the remaining non-flag arguments (operations)
	operations := flag.Args()

	// Check if operations are provided
	if len(operations) == 0 {
		fmt.Println("No operations provided")
		os.Exit(1)
	}

	// Print the operations
	fmt.Println("Operations:", strings.Join(operations, ", "))
}

func buildBloomFilterHash() error {
	bloomFilter := filter.NewBloomFilter(10000000, 5)
	file, err := os.Open("dict.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		bloomFilter.Add(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := dump(*bloomFilter); err != nil {
		return err
	}
	return nil
}

func dump(data any) error {
	f, err := os.Create("filtered.bin")
	if err != nil {
		return err
	}
	defer f.Close()
	err = binary.Write(f, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}
