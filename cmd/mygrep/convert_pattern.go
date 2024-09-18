package main

import (
	"fmt"
	"strings"
)

func convertPattern(regexp string) ([]RegExp, error) {
	var re []RegExp

	for i := 0; i < len(regexp); i++ {
		if i == 0 && regexp[i] == '^' {
			re = append(re, RegExp{Type: StartOfLine})
		} else if i == len(regexp)-1 && regexp[i] == '$' {
			re = append(re, RegExp{Type: EndOfLine})
		} else if regexp[i] == '\\' {
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
			token := RegExp{Type: Char, Char: rune(regexp[i])}
			// Look ahead for quantifiers
			if i+1 < len(regexp) && regexp[i+1] == '+' {
				token = RegExp{Type: Plus, Char: rune(regexp[i])}
				i++
			}
			re = append(re, token)

		}
	}
	return re, nil
}
