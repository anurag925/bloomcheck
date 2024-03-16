package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/anurag925/bloomcheck/src/filter"
)

const binFileName = "filter.bin"

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

	var filter filter.BloomFilter
	if err := loadFilter(&filter); err != nil {
		fmt.Println("error loading filter: ", err)
		return
	}

	fmt.Println("Words not found are: ")
	for _, e := range operations {
		fmt.Println(e, "found", test(filter, e))
	}

	// Print the operations
	fmt.Println("Words:", strings.Join(operations, ", "))
}

func buildBloomFilterHash() error {
	bloomFilter := filter.NewBloomFilter(1000000, 3)
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
	if err := bloomFilter.WriteToFile(binFileName); err != nil {
		return err
	}
	return nil
}

func loadFilter(filter *filter.BloomFilter) error {
	// Open the file for reading
	file, err := os.Open(binFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob decoder
	decoder := gob.NewDecoder(file)

	// Decode the BloomFilter struct from the file
	err = decoder.Decode(filter)
	if err != nil {
		return err
	}

	return nil
}

func test(filter filter.BloomFilter, word string) bool {
	return filter.Test([]byte(word))
}
