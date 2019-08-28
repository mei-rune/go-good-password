package good_password

import (
	"strings"
)

type negativeTest int

const (
	short negativeTest = iota

	// Negative tests that can be found multiple times in the same password.
	veryCommon
	commonWithSubsitution
	repeatedChar3
	repeatedPattern
	sequence3
)

func (n negativeTest) String() string {
	r := "<unknown>"
	switch n {
	case short:
		r = "short"
	case veryCommon:
		r = "contains very common word"
	case commonWithSubsitution:
		r = "contains very common word (written in leet speak)"
	case repeatedChar3:
		r = "repeated 3 or more characters"
	case repeatedPattern:
		r = "repeated pattern (3 characters or more)"
	case sequence3:
		r = "sequence of 3 or more characters"
	}
	return "-" + r
}

func (n negativeTest) Score() Score {
	switch n {
	// Too short is a huge penalty, nothing can redeem.
	case short:
		return -10
	case veryCommon, commonWithSubsitution:
		return -2
	default:
		return -1
	}
}

func negativeTests(pw string, extraWords []string) testResults {
	results := testResults{}

	// short, counting just runes (doesn't quite match good's behaviour, see
	// comment there, but "good" enough).
	l := 0
	for range pw {
		l++
	}
	if l < 8 {
		results = append(results, short)
	}

	// veryCommon + commonWithSubsitution
	wordList := append(append([]string{}, commonWordList...), extraWords...)
	lowerPw := strings.ToLower(pw)
	subPw := subsitutePassword(lowerPw)
	for _, word := range wordList {
		if strings.Contains(lowerPw, word) {
			results = append(results, veryCommon)
		} else if subPw != lowerPw && strings.Contains(subPw, word) {
			results = append(results, commonWithSubsitution)
		}
	}

	results = negativeSequences(pw, results)

	// repeatedPattern
	// Just ignore for very long passwords, to reduce chance of complexity attacks here.
	runes := []rune(pw)
	if len(runes) < 30 {
		for i := 0; i < len(runes)-3; i++ {
			part := runes[i : i+3]
			for j := i + 3; j <= len(runes)-3; j++ {
				if i != j && runeEq(part, runes[j:j+3]) {
					results = append(results, repeatedPattern)
				}
			}
		}
	}

	return results
}

func runeEq(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func negativeSequences(pw string, results testResults) testResults {
	var lastC rune
	repeats := 0
	seqF := 0
	seqR := 0
	for _, c := range pw {
		// repeatedChar3
		if lastC == c {
			repeats++
			if repeats == 2 {
				repeats = 0
				results = append(results, repeatedChar3)
			}
		} else {
			repeats = 0
		}

		// sequence3 (forward)
		if lastC+1 == c {
			seqF++
			if seqF == 2 {
				seqF = 0
				results = append(results, sequence3)
			}
		} else {
			seqF = 0
		}

		// sequence3 (back)
		if lastC-1 == c {
			seqR++
			if seqR == 2 {
				seqR = 0
				results = append(results, sequence3)
			}
		} else {
			seqR = 0
		}

		lastC = c
	}
	return results
}

func subsitutePassword(pw string) string {
	r := []rune{}
	for _, c := range pw {
		switch c {
		case '0':
			c = 'o'
		case '1', '!':
			c = 'i'
		case '3':
			c = 'e'
		case '4', '@':
			c = 'a'
		case '5', '$':
			c = 's'
		case '8':
			c = 'b'
		}
		r = append(r, c)
	}
	return string(r)
}
