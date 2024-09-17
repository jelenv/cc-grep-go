package main

import "fmt"

type REType int

const (
	Char REType = iota
	Digit
	Word
	CharGroup
	StartOfLine
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
	case StartOfLine:
		return "^"
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
