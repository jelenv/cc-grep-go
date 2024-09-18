package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Usage: echo <input_text> | mygrep -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	regexp := os.Args[2]

	inputText, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := match(regexp, inputText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("Matched: ", ok)

	if !ok {
		os.Exit(1)
	}
}

func match(regexp string, inputText []byte) (bool, error) {
	var result bool
	re, err := convertPattern(regexp)
	if err != nil {
		return false, err
	}
	if len(re) == 0 {
		return false, fmt.Errorf("empty pattern")
	}
	if len(inputText) == 0 {
		return false, fmt.Errorf("empty input text")
	}

	printRegExpArray(re)

	startOfLine := len(re) > 0 && re[0].Type == StartOfLine
	endOfLine := len(re) > 0 && re[len(re)-1].Type == EndOfLine

	if startOfLine {
		re = re[1:]
	}
	if endOfLine {
		re = re[:len(re)-1]
	}

	for start := 0; start < len(inputText); start++ {
		if startOfLine && start > 0 {
			break
		}

		match := true
		inputPos := start
		for i := 0; i < len(re); i++ {
			token := re[i]

			if token.Type == Plus {
				if inputPos >= len(inputText) || !matchToken(token, rune(inputText[inputPos])) {
					match = false
					break
				}
				for inputPos < len(inputText) && matchToken(token, rune(inputText[inputPos])) {
					inputPos++
				}
			} else if token.Type == QuestionMark {
				if inputPos < len(inputText) && matchToken(token, rune(inputText[inputPos])) {
					inputPos++
				}
				// continue without breaking even if the token doesn't match
			} else {
				if inputPos >= len(inputText) || !matchToken(token, rune(inputText[inputPos])) {
					// if negative char group did not match, then return immediately
					if token.Type == CharGroup && token.Negated {
						return false, nil
					}
					match = false
					break
				}
				inputPos++
			}
		}

		// check if we need to match end of line
		if match && endOfLine && start+len(re) != len(inputText) {
			match = false
		}

		if match {
			result = true
			break
		}
	}

	return result, nil
}

func matchToken(token RegExp, inputChar rune) bool {
	switch token.Type {
	case Char:
		return token.Char == inputChar
	case Digit:
		return unicode.IsDigit(inputChar)
	case Word:
		return unicode.IsLetter(inputChar) || unicode.IsDigit(inputChar) || inputChar == '_'
	case CharGroup:
		return strings.ContainsRune(string(token.CharArr), inputChar) != token.Negated
	case Plus, QuestionMark:
		return inputChar == token.Char
	}
	return false
}
