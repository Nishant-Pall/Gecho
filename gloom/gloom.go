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
	gloomArr []uint64
	seed     maphash.Seed
	len      uint64
	hash     maphash.Hash
	hashArr  []func(string) uint64
	hashLen  int
}

func (gloomFilter *GloomFilter) CreateGloomFilter(length uint64, hashes uint64, hashFunc GloomFilterHashFunc) error {
	if length < 1 {
		return fmt.Errorf("length cannot be less than 1")
	}

	gloomFilter.len = length
	gloomFilter.CreateGloomArr()
	gloomFilter.CreateSeed()
	gloomFilter.GenerateHashFunctions(hashes, hashFunc)
	gloomFilter.hashLen = len(gloomFilter.hashArr)

	return nil
}

func (f *GloomFilter) CreateGloomArr() {
	f.gloomArr = make([]uint64, f.len)
}

func (f *GloomFilter) CreateSeed() {
	f.seed = maphash.MakeSeed()
}

func (f *GloomFilter) GenerateHashFunctions(hashes uint64, hashFunc GloomFilterHashFunc) {
	f.hashArr = make([]func(string) uint64, hashes)

	for index := range f.hashArr {

		// generation should be outside invokation obviously
		n := rand.Uint64N(100)
		f.hashArr[index] = func(s string) uint64 {
			return MapHash(f, s) * n
		}
	}

}

func (f *GloomFilter) AddItem(s string) error {

	for _, hashFunc := range f.hashArr {
		hashInd := f.ModHash(hashFunc(s))

		f.gloomArr[hashInd] += 1
	}
	fmt.Printf("%v", f.gloomArr)
	return nil
}

func (f *GloomFilter) RemoveItem(s string) error {

	ok, _ := f.Lookup(s)
	if !ok {
		return fmt.Errorf("Key does not exist")
	}

	for _, hashFunc := range f.hashArr {
		hashInd := f.ModHash(hashFunc(s))

		if f.gloomArr[hashInd] > 0 {
			f.gloomArr[hashInd] -= 1
		}
	}
	return nil
}

func (f *GloomFilter) Lookup(s string) (bool, error) {

	for _, hashFunc := range f.hashArr {
		hashInd := f.ModHash(hashFunc(s))

		if f.gloomArr[hashInd] > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (f *GloomFilter) ModHash(hash uint64) uint64 {
	return hash % uint64(f.len)
}
