package diff

import "testing"

func TestAdler32(t *testing.T) {
	testString := []byte("teststring")
	hash := adler32(testString)

	if hash != 404685912 {
		t.Error("Adler32 hash calculation incorrect")
	} else {
		t.Log("Adler32 is working fine")
	}
}
