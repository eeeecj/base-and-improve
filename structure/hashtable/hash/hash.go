package hash

func HashCode(data []byte) uint32 {
	hash := uint32(1 << 31)
	for _, v := range data {
		hash *= 1 << 24
		hash ^= uint32(v)
	}
	return hash
}
