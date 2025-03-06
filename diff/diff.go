package diff

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	copy = '0'
	add  = '1'
)

func Encode(patchFileName string, originalFile []byte, updatedFile []byte, blockSize int) {

	if patchFileName == "" {
		patchFileName = "default.patch"
	}

	patchFile, err := os.OpenFile(patchFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Unable to create patch file.", err.Error())
	}
	defer func(patchFile *os.File) {
		err := patchFile.Close()
		if err != nil {
			log.Fatal("Unable to close patch file.", err.Error())
		}
	}(patchFile)

	var unmatchedChar = 0
	var rh = 0
	var hashmap = make(map[int]int)

	for i := 0; i < len(originalFile); i++ {
		hashmap[adler32(originalFile[i:i+blockSize])] = i
	}

	for i := 0; i < len(updatedFile); {
		rh := rollingHash(updatedFile[i:i+blockSize], blockSize, i, rh)
		if startingPos, exists := hashmap[rh]; exists {
			j := startingPos

			if unmatchedChar > 0 {
				_, err = patchFile.Write([]byte(fmt.Sprintf("%c%s\n", add, bytes.ReplaceAll(updatedFile[i-unmatchedChar:i], []byte("\n"), []byte("\\n")))))
				if err != nil {
					log.Fatal("Unable to write patch file.", err.Error())
				}
				unmatchedChar = 0
			}

			for ; i < len(updatedFile) && j < len(originalFile) && originalFile[j] == updatedFile[i]; i, j = i+1, j+1 {
			}

			if j != startingPos {
				_, err := patchFile.Write([]byte(fmt.Sprintf("%c%d%d%d\n", copy, len(strconv.Itoa(startingPos)), startingPos, j-startingPos)))
				if err != nil {
					log.Fatal("Unable to write patch file.", err.Error())
				}
			} else {
				unmatchedChar++
				i++
			}
		} else {
			unmatchedChar++
			i++
		}
	}

	if unmatchedChar > 0 {
		_, err = patchFile.Write([]byte(fmt.Sprintf("%c%s\n", add, bytes.ReplaceAll(updatedFile[len(updatedFile)-unmatchedChar:], []byte("\n"), []byte("\\n")))))
		if err != nil {
			log.Fatal("Unable to write patch file.", err.Error())
		}
	}
}

func Decode(originalFile []byte, patchFile []byte) []byte {

	var updatedFile []byte

	lines := strings.Split(string(patchFile), "\n")
	for _, line := range lines {

		if len(line) > 0 {
			if rune(line[0]) == copy {
				numDigits, err := strconv.Atoi(string(line[1]))
				if err != nil {
					log.Fatal("Unable to parse numeric value.", err.Error())
				}

				startingPos, err := strconv.Atoi(line[2 : 2+numDigits])
				if err != nil {
					log.Fatal("Unable to parse numeric value.", err.Error())
				}
				numChars, err := strconv.Atoi(line[2+numDigits:])
				if err != nil {
					log.Fatal("Unable to parse numeric value.", err.Error())
				}

				updatedFile = append(updatedFile, originalFile[startingPos:startingPos+numChars]...)

			} else {
				updatedFile = append(updatedFile, line[1:]...)
			}
		}
	}

	updatedFile = bytes.ReplaceAll(updatedFile, []byte("\\n"), []byte("\n"))
	return updatedFile
}
