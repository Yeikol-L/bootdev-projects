package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		os.Exit(1)
	}
	currentLine := ""
	for {
		chunck := make([]byte, 8)
		_, err := file.Read(chunck)
		parts := strings.Split(string(chunck), "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Println("read: %")
		}
		if err != nil {
			os.Exit(0)
		}

		// fmt.Printf("read: %s\n", chunck)
	}
}
