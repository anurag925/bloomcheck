package filter

import (
	"hash/fnv"
)

// BloomFilter is a simple implementation of a Bloom filter.
type BloomFilter struct {
	bitset []bool
	m      uint
	k      uint
	n      uint
}

// NewBloomFilter creates a new Bloom filter with the given size and number of hash functions.
func NewBloomFilter(m, k uint) *BloomFilter {
	return &BloomFilter{
		bitset: make([]bool, m),
		m:      m,
		k:      k,
		n:      0,
	}
}

// Add adds an element to the Bloom filter.
func (bf *BloomFilter) Add(data []byte) {
	for i := uint(0); i < bf.k; i++ {
		hash := bf.hash(data, i)
		bf.bitset[hash%bf.m] = true
	}
	bf.n++
}

// Test checks if an element is in the Bloom filter.
func (bf *BloomFilter) Test(data []byte) bool {
	for i := uint(0); i < bf.k; i++ {
		hash := bf.hash(data, i)
		if !bf.bitset[hash%bf.m] {
			return false
		}
	}
	return true
}

// hash generates a hash value for the given data and seed.
func (bf *BloomFilter) hash(data []byte, seed uint) uint {
	hash := fnv.New64a()
	hash.Write(data)
	hash.Write([]byte{byte(seed)})
	return uint(hash.Sum64()) % bf.m
}
