package filter

import (
	"encoding/gob"
	"hash/fnv"
	"os"
)

// BloomFilter is a simple implementation of a Bloom filter.
type BloomFilter struct {
	BitSet []bool
	// size of the dataset
	M uint64
	// number of hashes being done
	K uint64
	// length of the bit set
	N uint64
}

// NewBloomFilter creates a new Bloom filter with the given size and number of hash functions.
func NewBloomFilter(m, k uint64) *BloomFilter {
	return &BloomFilter{
		BitSet: make([]bool, m),
		M:      m,
		K:      k,
		N:      0,
	}
}

// Add adds an element to the Bloom filter.
func (bf *BloomFilter) Add(data []byte) {
	for i := uint64(0); i < bf.K; i++ {
		hash := bf.hash(data, i)
		bf.BitSet[hash%bf.M] = true
	}
	bf.N++
}

// Test checks if an element is in the Bloom filter.
func (bf *BloomFilter) Test(data []byte) bool {
	for i := uint64(0); i < bf.K; i++ {
		hash := bf.hash(data, i)
		if !bf.BitSet[hash%bf.M] {
			return false
		}
	}
	return true
}

// hash generates a hash value for the given data and seed.
func (bf *BloomFilter) hash(data []byte, seed uint64) uint64 {
	hash := fnv.New64a()
	hash.Write(data)
	hash.Write([]byte{byte(seed)})
	return uint64(hash.Sum64()) % bf.M
}

func (bf *BloomFilter) WriteToFile(filename string) error {
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob encoder
	encoder := gob.NewEncoder(file)

	// Encode and write the BloomFilter struct
	err = encoder.Encode(bf)
	if err != nil {
		return err
	}

	return nil
}
