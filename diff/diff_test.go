package diff

import (
	"io"
	"log"
	"os"
	"testing"
)

func fileOpenErrorHandler(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s. %e", fileName, err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error when closing the file %s. %e", fileName, err)
		}
	}(file)

	filebytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file %s. %e", fileName, err)
	}

	return filebytes
}

func TestEncode(t *testing.T) {
	var f1 = "../testfiles/file1.txt"
	var f2 = "../testfiles/file2.txt"
	var patchFile = "../testfiles/test1.patch"
	var requiredResult = "01045\n1This line makes this file different\n024238\n1\\n\\nTest Case: Bufio reader\n"

	Encode(patchFile, fileOpenErrorHandler(f1), fileOpenErrorHandler(f2), 8)

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

func TestEncodeOpposite(t *testing.T) {
	var f1 = "../testfiles/file2.txt"
	var f2 = "../testfiles/file1.txt"
	var patchFile = "../testfiles/test2.patch"
	var requiredResult = "01045\n028335\n"

	Encode(patchFile, fileOpenErrorHandler(f1), fileOpenErrorHandler(f2), 8)

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
	var patchFile = "../testfiles/test1.patch"

	requiredResult, err := os.ReadFile(f2)
	if err != nil {
		t.Error(err)
	}

	actualResult := Decode(fileOpenErrorHandler(f1), fileOpenErrorHandler(patchFile))

	if string(actualResult) != string(requiredResult) {
		t.Error("Actual result does not match required result. Actual result: ", string(actualResult), "Required result: ", string(requiredResult))
	} else {
		t.Log("Actual result matches with the required result")
	}
}

func TestDecodeOpposite(t *testing.T) {
	var f1 = "../testfiles/file2.txt"
	var f2 = "../testfiles/file1.txt"
	var patchFile = "../testfiles/test2.patch"

	requiredResult, err := os.ReadFile(f2)
	if err != nil {
		t.Error(err)
	}

	actualResult := Decode(fileOpenErrorHandler(f1), fileOpenErrorHandler(patchFile))

	if string(actualResult) != string(requiredResult) {
		t.Error("Actual result does not match required result. Actual result: ", string(actualResult), "Required result: ", string(requiredResult))
	} else {
		t.Log("Actual result matches with the required result")
	}
}
