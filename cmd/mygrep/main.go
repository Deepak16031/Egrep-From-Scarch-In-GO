package main

import (
	"bytes"
	"io"
	"strings"

	// Uncomment this to pass the first stage
	// "bytes"
	"fmt"
	"os"
)

const digitClass = "\\d"
const alphaNumericClass = "\\w"

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	//if len(os.Args) < 3 || os.Args[1] != "-E" {
	//	fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
	//	os.Exit(2) // 1 means no lines were selected, >1 means error
	//}

	pattern := os.Args[2]
	//pattern := "cat$"
	line, _ := io.ReadAll(os.Stdin)
	//line := []byte("cat")

	ok := match(line, pattern)

	//ok, err := matchLine(line, pattern)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	//	os.Exit(2)
	//}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

//func exitOnOk(ok bool) {
//	if !ok {
//		os.Exit(1)
//	} else {
//		os.Exit(0)
//	}
//}

func match(line []byte, pattern string) bool {
	//try to match first char
	if pattern[0] == '^' {
		return matchUtil(line, pattern[1:])
	} else if pattern[0] == '\\' {
		if pattern[1] == 'd' {
			for {
				ok, indx := containsDigit(line)
				if !ok {
					return false
				}
				if matchUtil(line[indx+1:], pattern[2:]) {
					return true
				}
				line = line[indx+1:]
			}
		} else if pattern[1] == 'w' {
			for {
				ok, indx := isAlphaNumeric(line)
				if !ok {
					return false
				}
				if matchUtil(line[indx+1:], pattern[2:]) {
					return true
				}
				line = line[indx+1:]
			}
		}
		return false
	} else if pattern[0] == '[' {
		allChars := pattern[1 : strings.IndexByte(pattern[1:], ']')+1]
		classChars := allChars
		notCheck := pattern[1] == '^'
		if notCheck {
			classChars = classChars[1:]
		}
		for {
			indx := bytes.IndexAny(line, classChars)
			result := func() bool {
				if notCheck {
					return indx == -1
				}
				return indx >= 0
			}()
			if !result {
				return false
			}
			if matchUtil(line[indx+1:], pattern[len(allChars)+2:]) {
				return true
			}
			line = line[indx+1:]
		}
	}
	return matchUtil(line, pattern)

}

// this is used to match exact characters with pattern
func matchUtil(line []byte, pattern string) bool {
	lineIndx := 0
	i := 0
	for i = 0; i < len(pattern) && lineIndx < len(line); i++ {
		r := pattern[i]
		if r == '\\' {
			if pattern[i+1] == 'd' {
				ok, indx := containsDigit(line[lineIndx:])
				if !ok || indx != 0 {
					return false
				}
			} else if pattern[i+1] == 'w' {
				ok, indx := isAlphaNumeric(line[lineIndx:])
				if !ok || indx != 0 {
					return false
				}
			}
			lineIndx++
			i += 1
		} else if r == '[' {
			classEndIndx := strings.IndexByte(pattern[i+1:], ']')
			allChars := pattern[i+1 : i+1+classEndIndx]
			notCheck := pattern[i+1] == '^'
			if notCheck {
				allChars = allChars[1:]
			}
			indx := strings.IndexAny(allChars, string(line[lineIndx]))
			result := func() bool {
				if notCheck {
					return indx == -1
				}
				return indx >= 0
			}()
			if !result {
				return false
			}
			lineIndx++
			i += classEndIndx + 1
		} else if r == '$' {
			return len(line) == lineIndx
		} else if r == line[lineIndx] {
			lineIndx++
		} else {
			return false
		}
	}
	return i == len(pattern) || (pattern[i] == '$') // for loop breaking conditions simplifies the logic
}

//func matchCharClassInStart(line []byte, pattern string) (bool, bool) {
//	if pattern[0:2] == digitClass {
//		ok, indx := containsDigit(line)
//		if !ok {
//			return false, true
//		}
//		return match(line, pattern, indx+1, 2), true
//	} else if pattern[0:2] == alphaNumericClass {
//		ok, indx := isAlphaNumeric(line)
//		if !ok {
//			return false, true
//		}
//		return match(line, pattern, indx+1, 2), true
//	}
//	return false, false
//}

func startsWithCharClass(pattern string, patternIndx int) bool {
	return len(pattern) > 1 && patternIndx == 0 && pattern[0] == '\\'
}

func containsDigit(line []byte) (bool, int) {
	digits := "0123456789"
	var ok int
	ok = bytes.IndexAny(line, digits)
	return ok >= 0, ok
}

func isAlphaNumeric(line []byte) (bool, int) {
	for i, r := range line {
		asciiValue := int(r)
		isDigit := asciiValue >= 48 && asciiValue <= 57
		isLiteral := (asciiValue >= 65 && asciiValue <= 90) || (asciiValue >= 97 && asciiValue <= 122)
		if isDigit || isLiteral {
			return true, i
		}
	}
	return false, -1
}

func matchNotExist(line []byte, pattern string) (bool, int) {
	ok, indx := matchLine(line, pattern)
	return !ok, indx
}

func matchLine(line []byte, pattern string) (bool, int) {
	//if utf8.RuneCountInString(pattern) != 1 {
	//	return false, fmt.Errorf("unsupported pattern: %q", pattern)
	//}

	var ok int

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this to pass the first stage
	ok = bytes.IndexAny(line, pattern)

	//return ok, nil
	return ok >= 0, ok
}
