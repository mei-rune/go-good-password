// Package good_password tests a password's strength.
package good_password

import (
	"fmt"
	"regexp"
)

var (
	notWordRE = regexp.MustCompile(`[^\p{L}\p{N}]`)
)

type scorer interface {
	fmt.Stringer
	Score() Score
}

type testResults []scorer

type Score int

func (s Score) String() string {
	switch {
	case s < 1:
		return "Terrible"
	case s <= 2:
		return "Weak"
	case s <= 3:
		return "Okay"
	case s <= 4:
		return "Good"
	default:
		return "Strong"
	}
}

// Check runs checks against the given password and returns a Score.
// extraWords specify additional words to treat as "common" (e.g. the user's
// name), it should be a lowercase list of short words.
// The list of strings returned contains notes about the password, the first
// character is "-" or "+", the rest a comment.
func Check(password string, extraWords []string) (Score, []string) {
	results := append(append(testResults{}, goodTests(password)...), negativeTests(password, extraWords)...)
	score := Score(0)
	info := make(map[scorer]bool)
	for _, result := range results {
		score += result.Score()
		info[result] = true
	}
	var infoList []string
	for r := range info {
		infoList = append(infoList, r.String())
	}
	return score, infoList
}

// ExtractWords takes words out of the given strings. It is unicode aware.
// The output of this can be passed directly to Check's extraWord parameter.
func ExtractWords(words ...string) []string {
	results := []string{}
	for _, word := range words {
		for _, f := range notWordRE.Split(word, -1) {
			if len(f) >= 3 {
				results = append(results, f)
			}
		}
	}
	return results
}
