package good_password_test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mei-rune/go-good-password"
)

func ExampleCheck() {
	// Call Check, it's recommended to include the user's name to discourage passwords containing their name.
	score, info := good_password.Check("good-password?", []string{"username"})

	// Only needed to make "Output" below consistent.
	sort.Strings(info)

	// Score has a String here, but you may also wish to check: score >= good_password.RecommendedScore
	fmt.Printf("%v password! (%s)\n", score, strings.Join(info, ", "))

	// Output: Terrible password! (+has lowercase letter, +has special character, +length at least 13 characters, -contains very common word)
}
