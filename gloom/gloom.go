package gloom

import (
	"fmt"
	"hash/fnv"
	"hash/maphash"
	"math/rand/v2"
)

func NewGloomFilter() *GloomFilter {
	return new(GloomFilter)
}

func MapHash(f *GloomFilter, s string) uint64 {

	f.hash.SetSeed(f.seed)
	f.hash.WriteString(s)
	str := f.hash.Sum64()

	return str
}

func BasicHash(f *GloomFilter, s string) uint64 {

	h := fnv.New32a()
	h.Write([]byte(s))

	return uint64(h.Sum32())
}

type GloomFilterHashFunc func(*GloomFilter, string) uint64

type GloomFilter struct {
	gloomArr []int
	seed     maphash.Seed
	len      int
	hash     maphash.Hash
	hashes   []func(string) uint64
	hashLen  int
}

func (gloomFilter *GloomFilter) InstantiateGloomFilter(length int, hashes int, hashFunc GloomFilterHashFunc) error {
	if length < 1 {
		return fmt.Errorf("length cannot be less than 1")
	}

	gloomFilter.len = length
	gloomFilter.CreateGloomArr()
	gloomFilter.CreateSeed()
	gloomFilter.GenerateHashFunctions(hashes, hashFunc)
	gloomFilter.hashLen = len(gloomFilter.hashes)

	return nil
}

func (f *GloomFilter) CreateGloomArr() {
	f.gloomArr = make([]int, f.len)
}

func (f *GloomFilter) CreateSeed() {
	f.seed = maphash.MakeSeed()
}

func (f *GloomFilter) GenerateHashFunctions(hashes int, hashFunc GloomFilterHashFunc) {
	f.hashes = make([]func(string) uint64, hashes)

	for index := range f.hashes {

		// generation should be outside invokation obviously
		n := rand.Uint64N(100)
		f.hashes[index] = func(s string) uint64 {
			return MapHash(f, s) * n
		}
	}

}

func (f *GloomFilter) AddItem(s string) {

	for _, hashFunc := range f.hashes {
		hashInd := f.ModHash(hashFunc(s))

		f.gloomArr[hashInd] = 1
	}
}

func (f *GloomFilter) CheckMembership(s string) bool {

	for _, hashFunc := range f.hashes {
		hashInd := f.ModHash(hashFunc(s))

		if f.gloomArr[hashInd] != 1 {
			return false
		}
	}
	return true
}

func (f *GloomFilter) ModHash(hash uint64) uint64 {
	return hash % uint64(f.len)
}

// func test() {
// 	len := 40
// 	hashes := 5

// 	gloomFilter := gloom.NewGloomFilter()
// 	gloomFilter.InstantiateGloomFilter(len, hashes, gloom.MapHash)

// 	gloomFilter.AddItem("NISHANT")
// 	gloomFilter.AddItem("ARUSHI")

// 	fmt.Println(gloomFilter.CheckMembership("NISHANT"))
// 	fmt.Println(gloomFilter.CheckMembership("ARUSH1"))

//
