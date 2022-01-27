package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	wordFile *os.File
	fileInfo os.FileInfo
	fileSize int64
	err      error

	word []byte
)

type coloredChar struct {
	char  byte
	color *color.Color
}

func main() {
	wordFile, err = os.Open("assets/words.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer wordFile.Close()

	fileInfo, err = os.Stat("assets/words.txt")

	if err != nil {
		log.Fatal(err)
	}

	fileSize = fileInfo.Size()

	word = getWord()

	fmt.Printf("word: %s\n", word)

	guessWord := []byte("adieu")

	response := guess(guessWord)

	if response != nil {
		for _, c := range response {
			c.color.Printf("%s", string(c.char))
		}
		//log.Println()
		fmt.Println()
	} else {
		log.Printf("response null")
	}
}

func getWord() []byte {
	var numLines int = int(fileSize / 6)

	rand.Seed(time.Now().UnixNano())

	var randLine int = rand.Intn(numLines + 1)
	var offset int64 = int64(randLine * 6)

	_, err = wordFile.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	word := make([]byte, 5)
	_, err = wordFile.Read(word)

	if err != nil {
		log.Fatal(err)
	}

	return word
}

func guess(guess []byte) []*coloredChar {
	response := make([]*coloredChar, 5)
	for i, char := range guess {
		if char == word[i] {
			response[i] = &coloredChar{char, color.New(color.FgHiGreen)}
		} else if bytes.Contains(word, []byte{char}) {
			response[i] = &coloredChar{char, color.New(color.FgHiYellow)}
		} else {
			response[i] = &coloredChar{char, color.New(color.FgWhite)}
		}
	}

	return response
}
