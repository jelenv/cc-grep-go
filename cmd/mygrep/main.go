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

type REType int

const (
	Char REType = iota
	Digit
	Word
	CharGroup
)

type RegExp struct {
	Type    REType
	Char    rune
	CharArr []rune
	Negated bool
}

func (r RegExp) String() string {
	switch r.Type {
	case Char:
		if r.Char == '\\' {
			return fmt.Sprintf("\\%c", r.Char)
		}
		return string(r.Char)
	case Digit:
		return "\\d"
	case Word:
		return "\\w"
	case CharGroup:
		var negated string
		if r.Negated {
			negated = "^"
		}
		return fmt.Sprintf("[%s%s]", negated, string(r.CharArr))
	}
	return ""
}

func printRegExpArray(re []RegExp) {
	fmt.Printf("RegExp: ")
	for _, r := range re {
		fmt.Printf("%s", r.String())
	}
	fmt.Println()
}

func convertPattern(regexp string) ([]RegExp, error) {
	var re []RegExp

	for i := 0; i < len(regexp); i++ {
		if regexp[i] == '\\' {
			i++
			if i < len(regexp) {
				switch regexp[i] {
				case 'd':
					re = append(re, RegExp{Type: Digit})
				case 'w':
					re = append(re, RegExp{Type: Word})
				default:
					re = append(re, RegExp{Type: Char, Char: '\\'})
					re = append(re, RegExp{Type: Char, Char: rune(regexp[i])})
				}
			} else {
				return nil, fmt.Errorf("unsupported pattern: %q", regexp)
			}
		} else if regexp[i] == '[' {
			end := strings.IndexRune(regexp[i+1:], ']')
			if end == -1 {
				return nil, fmt.Errorf("invalid character group: %q", regexp)
			}
			negated := false

			if i+1 < len(regexp) && regexp[i+1] == '^' {
				negated = true
				i++
			}
			charGroup := regexp[i+1 : i+end]
			re = append(re, RegExp{Type: CharGroup, CharArr: []rune(charGroup), Negated: negated})
			i += end + 1
		} else {
			re = append(re, RegExp{Type: Char, Char: rune(regexp[i])})
		}
	}
	return re, nil
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

func match(regexp string, inputText []byte) (bool, error) {
	result := false

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

	for start := 0; start < len(inputText); start++ {
		match := true
		for i := 0; i < len(re); i++ {
			// fmt.Println("Token: ", re[i].String(), "Input: ", string(inputText[start+i]), "Match: ", matchToken(re[i], rune(inputText[start+i])))
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
