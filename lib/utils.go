package lib

import "hash/fnv"

func Hash(s string) uint32 {
	if s == "" {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
