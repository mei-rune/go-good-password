package good_password

import (
	"unicode"
)

type goodTest int

const (
	length10 goodTest = iota
	length13
	length16

	hasAlphaUpper
	hasAlphaLower
	hasDigit
	hasOther
)

func (g goodTest) String() string {
	r := "<unknown>"
	switch g {
	case length10:
		r = "length at least 10 characters"
	case length13:
		r = "length at least 13 characters"
	case length16:
		r = "length at least 16 characters"
	case hasAlphaUpper:
		r = "has uppercase letter"
	case hasAlphaLower:
		r = "has lowercase letter"
	case hasDigit:
		r = "has digit"
	case hasOther:
		r = "has special character"
	}
	return "+" + r
}

func (g goodTest) Score() Score {
	switch g {
	case length16:
		return 3
	case length13:
		return 2
	}
	return 1
}

func goodTests(pw string) testResults {
	results := make(map[goodTest]bool)

	l := 0
	for _, r := range pw {
		switch {
		case unicode.IsDigit(r):
			results[hasDigit] = true
		case unicode.IsUpper(r):
			results[hasAlphaUpper] = true
		case unicode.IsLower(r):
			results[hasAlphaLower] = true
		case unicode.IsGraphic(r):
			results[hasOther] = true
		default:
			l--
		}
		l++
	}

	// Length is based on counting runes that are graphics (may overcount with
	// combining characters, but better than counting an emoji password as very
	// good).
	if l >= 16 {
		results[length16] = true
	} else if l >= 13 {
		results[length13] = true
	} else if l >= 10 {
		results[length10] = true
	}

	rs := testResults{}
	for r := range results {
		rs = append(rs, r)
	}
	return rs
}
