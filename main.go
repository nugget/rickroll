package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func singLine(line string) error {
	fmt.Println("# LINE", line)

	fmt.Printf("\n")
	return nil
}

func main() {
	filename := "lyrics.dat"

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewBuffer(b))
	for {
		line, prefix, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(line, prefix, err)
		}

		if len(line) == 0 {
			fmt.Printf("\n")
		} else {
			err = singLine(string(line))
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
