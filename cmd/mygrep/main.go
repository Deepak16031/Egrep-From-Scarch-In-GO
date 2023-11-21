package main

import (
	"bytes"
	"io"

	// Uncomment this to pass the first stage
	// "bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	//line := []byte("alpha-num3ric") // assume we're only dealing with a single line
	line, err := io.ReadAll(os.Stdin)

	if "\\d" == pattern {
		ok := containsDigit(line)
		exitOnOk(ok)
	} else if "\\w" == pattern {
		ok := isAlphaNumbers(line)
		exitOnOk(ok)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func exitOnOk(ok bool) {
	if !ok {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func containsDigit(line []byte) bool {
	digits := "0123456789"
	var ok bool
	ok = bytes.ContainsAny(line, digits)
	return ok
}

func isAlphaNumbers(line []byte) bool {
	for _, r := range line {
		asciiValue := int(r)
		isDigit := asciiValue >= 48 && asciiValue <= 57
		isLiteral := (asciiValue >= 65 && asciiValue <= 90) || (asciiValue >= 97 && asciiValue <= 122)
		if isDigit || isLiteral {
			return true
		}
	}
	return false
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) != 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this to pass the first stage
	ok = bytes.ContainsAny(line, pattern)

	return ok, nil
}
