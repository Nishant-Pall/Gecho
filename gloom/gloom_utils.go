package gloom

import "hash/fnv"

type GloomFilterHashFunc func(*BaseGloomFilter, string) uint64

func NewGloomFilter() *BaseGloomFilter {
	return new(BaseGloomFilter)
}

func MapHash(f *BaseGloomFilter, s string) uint64 {

	f.hash.SetSeed(f.seed)
	f.hash.WriteString(s)
	str := f.hash.Sum64()

	return str
}

func BasicHash(f *BaseGloomFilter, s string) uint64 {

	h := fnv.New32a()
	h.Write([]byte(s))

	return uint64(h.Sum32())
}
