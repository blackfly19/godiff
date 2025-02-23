package diff

const prime = 65521

func rollingHash(data []byte, blockSize int, pos int, previousHash int) int {

	if previousHash == 0 {
		return adler32(data)
	}
	s1 := (previousHash - int(data[pos-blockSize]+data[pos])) % prime
	s2 := (previousHash - blockSize*int(data[pos-blockSize]) + s1) % prime
	return s2<<16 + s1
}

func adler32(data []byte) int {
	var s1 = 1
	var s2 = 0
	var prime = 65521

	for i := 0; i < len(data); i++ {
		s1 = (s1 + int(data[i])) % prime
		s2 = (s2 + s1) % prime
	}
	return s2<<16 + s1
}
