package gloom

import "hash/fnv"

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
