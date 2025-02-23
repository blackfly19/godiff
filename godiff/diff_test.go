package godiff

import (
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	var f1 = "../testfiles/file1.txt"
	var f2 = "../testfiles/file2.txt"
	var patchFile = "../testfiles/test.patch"
	var requiredResult = "01045\n1This line makes this file different\n024238\n1\\n\\nTest Case: Bufio reader\n"

	Encode(patchFile, f1, f2, 8)

	actualResult, err := os.ReadFile(patchFile)
	if err != nil {
		t.Error(err)
	}

	if string(actualResult) != requiredResult {
		t.Error("Actual result does not match required result. Actual result: ", string(actualResult), "Required result: ", requiredResult)
	} else {
		t.Log("Actual result matches with the required result")
	}

}

func TestDecode(t *testing.T) {
	var f1 = "../testfiles/file1.txt"
	var f2 = "../testfiles/file2.txt"
	var patchFile = "../testfiles/test.patch"

	requiredResult, err := os.ReadFile(f2)
	if err != nil {
		t.Error(err)
	}

	actualResult := Decode(f1, patchFile)

	if string(actualResult) != string(requiredResult) {
		t.Error("Actual result does not match required result. Actual result: ", string(actualResult), "Required result: ", string(requiredResult))
	} else {
		t.Log("Actual result matches with the required result")
	}
}
