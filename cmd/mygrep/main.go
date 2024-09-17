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
	if startOfLine {
		re = re[1:]
	}

	for start := 0; start < len(inputText); start++ {
		if startOfLine && start > 0 {
			break
		}

		match := true
		for i := 0; i < len(re); i++ {
			// if start+i < len(inputText) {
			// 	fmt.Println("Token: ", re[i].String(), "Input: ", string(inputText[start+i]), "Match: ", matchToken(re[i], rune(inputText[start+i])))
			// } else {
			// 	fmt.Println("Token: ", re[i].String(), "Input: ", "EOF", "Match: ", false)
			// }
			if start+i >= len(inputText) || !matchToken(re[i], rune(inputText[start+i])) {

				// if negative char group did not match, then return immediately
				if re[i].Type == CharGroup && re[i].Negated {
					fmt.Println("Matched: ", false)
					return false, nil
				}
				match = false
				break
			}
		}
		if match {
			result = true
			break
		}
	}

	fmt.Println("Matched: ", result)
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
	}
	return false
}
