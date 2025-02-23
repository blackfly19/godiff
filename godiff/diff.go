package godiff

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	COPY = '0'
	ADD  = '1'
)

func Encode(patchFileName string, originalFileName string, updatedFileName string, blockSize int) {

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

	originalFile, err := os.ReadFile(originalFileName)
	if err != nil {
		log.Fatal("Unable to read original file.", err.Error())
	}

	updatedFile, err := os.ReadFile(updatedFileName)
	if err != nil {
		log.Fatal("Unable to read updated file.", err.Error())
	}

	var unmatchedChar = 0
	var rh = 0
	var hashmap = make(map[int]int)

	for i := 0; i < len(originalFile); i++ {
		hashmap[adler32(originalFile[i:i+blockSize])] = i
	}

	for i := 0; i < len(updatedFile); i++ {
		rh := rollingHash(updatedFile[i:i+blockSize], blockSize, i, rh)
		if startingPos, exists := hashmap[rh]; exists {
			j := startingPos

			if unmatchedChar > 0 {
				_, err = patchFile.Write([]byte(fmt.Sprintf("%c%s\n", ADD, bytes.ReplaceAll(updatedFile[i-unmatchedChar-1:i], []byte("\n"), []byte("\\n")))))
				if err != nil {
					log.Fatal("Unable to write patch file.", err.Error())
				}
				unmatchedChar = 0
			}

			for ; i < len(updatedFile) && j < len(originalFile) && originalFile[j] == updatedFile[i]; j++ {
				i++
			}
			_, err := patchFile.Write([]byte(fmt.Sprintf("%c%d%d%d\n", COPY, len(strconv.Itoa(startingPos)), startingPos, j-startingPos)))
			if err != nil {
				log.Fatal("Unable to write patch file.", err.Error())
			}
		} else {
			unmatchedChar++
		}
	}

	if unmatchedChar > 0 {
		_, err = patchFile.Write([]byte(fmt.Sprintf("%c%s\n", ADD, bytes.ReplaceAll(updatedFile[len(updatedFile)-unmatchedChar-1:], []byte("\n"), []byte("\\n")))))
		if err != nil {
			log.Fatal("Unable to write patch file.", err.Error())
		}
	}
}

func Decode(originalFileName string, patchFileName string) []byte {

	originalFile, err := os.ReadFile(originalFileName)
	if err != nil {
		log.Fatal("Unable to read original file.", err.Error())
	}

	patchFile, err := os.Open(patchFileName)
	if err != nil {
		log.Fatal("Unable to open patch file.", err.Error())
	}

	var updatedFile []byte
	scanner := bufio.NewScanner(patchFile)
	for scanner.Scan() {
		line := scanner.Text()

		if rune(line[0]) == COPY {
			numDigits, err := strconv.Atoi(string(line[1]))
			if err != nil {
			}

			startingPos, err := strconv.Atoi(line[2 : 2+numDigits])
			fmt.Println(startingPos)
			if err != nil {
			}
			numChars, err := strconv.Atoi(line[2+numDigits:])
			if err != nil {
			}

			updatedFile = append(updatedFile, originalFile[startingPos:startingPos+numChars]...)

		} else {
			updatedFile = append(updatedFile, line[1:]...)
		}
	}

	updatedFile = bytes.ReplaceAll(updatedFile, []byte("\\n"), []byte("\n"))
	return updatedFile
}
