package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	fmt.Printf("File size: %d", fileSize)
	word = getWord()
	fmt.Printf("Word length: %d\n", len(word))
	fmt.Printf("word: %s\n", word)

	for {
		attempts := 0
		//guessWord := []byte("adieu")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Guess: ")
		text, _ := reader.ReadBytes('\n')

		response, correct := guess(text)

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

}

func getWord() []byte {
	rand.Seed(time.Now().UnixNano())
	numLines, err := lineCounter(wordFile)
	fmt.Printf("num lines: %d\n", numLines)
	var lineSize int64 = fileSize / int64(numLines)
	if err != nil {
		log.Fatal("panic")
	}

	var randLine int64 = int64(rand.Intn(int(numLines) + 1))
	fmt.Printf("rand line: %d\n", randLine)
	var offset int64 = int64(randLine * lineSize)

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

func guess(guess []byte) ([]*coloredChar, bool) {
	response := make([]*coloredChar, 5)

	var correct int = 0

	for i, char := range guess {
		if char == word[i] {
			correct++
			response[i] = &coloredChar{char, color.New(color.FgHiGreen)}
		} else if bytes.Contains(word, []byte{char}) {
			response[i] = &coloredChar{char, color.New(color.FgHiYellow)}
		} else {
			response[i] = &coloredChar{char, color.New(color.FgWhite)}
		}
	}
	if correct == 5 {
		return response, true
	}

	return response, false
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
