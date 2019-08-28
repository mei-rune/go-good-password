package good_password

import (
	"reflect"
	"sort"
	"testing"
)

func TestCheck(t *testing.T) {
	for _, test := range []struct {
		password      string
		extraWords    []string
		expectedScore Score
		expectedInfo  []string
	}{
		{"", nil, -10, []string{"-short"}},
		{"foobarbb", nil, 1, []string{"+has lowercase letter"}},
		// Words
		{"goodPassword", nil, -1, []string{"+has lowercase letter", "+has uppercase letter", "+length at least 10 characters", "-contains very common word"}},
		{"betterP@ssword", nil, 1, []string{"+has lowercase letter", "+has special character", "+has uppercase letter", "+length at least 13 characters", "-contains very common word (written in leet speak)"}},
		{"greatP@ssw()rd", nil, 3, nil},
		{"wordUPdude!", nil, 4, nil},
		{"foobarbb", []string{"foo"}, -1, nil},
		{"hunterbb", []string{"food"}, -1, nil},
		// Repeated character
		{"ooofooba", nil, 0, []string{"+has lowercase letter", "-repeated 3 or more characters"}},
		// Sequences
		{"abcfooba", nil, 0, []string{"+has lowercase letter", "-sequence of 3 or more characters"}},
		{"edcfooba", nil, 0, nil},
		{"123fooba", nil, 1, nil},
		{"654fooba", nil, 1, nil},
		// Repeated pattern
		{"fooblabla", nil, 0, []string{"+has lowercase letter", "-repeated pattern (3 characters or more)"}},
		// Not checked for >= 30
		{"blablablablablablablablablabla", nil, 4, []string{"+has lowercase letter", "+length at least 16 characters"}},
		// Unicode
		{"caféfooba", nil, 1, []string{"+has lowercase letter"}},
		{"cafÉfooba", nil, 2, []string{"+has lowercase letter", "+has uppercase letter"}},
		{"caf♜barbaz", nil, 3, []string{"+has lowercase letter", "+has special character", "+length at least 10 characters"}},
		// Still considered short, we count runes.
		{"🤓🤓🤓🤓🤓🤓🤓", nil, -15, []string{"+has special character", "-repeated 3 or more characters", "-repeated pattern (3 characters or more)", "-short"}},
		// Good password (yeah, not really)
		{"foo-9AR-baz", nil, 5, []string{"+has digit", "+has lowercase letter", "+has special character", "+has uppercase letter", "+length at least 10 characters"}},
		// Great password
		{"foo-9ARbazHAS", nil, 6, []string{"+has digit", "+has lowercase letter", "+has special character", "+has uppercase letter", "+length at least 13 characters"}},
		// Greatest password
		{"h.unter2something3THING", nil, 7, []string{"+has digit", "+has lowercase letter", "+has special character", "+has uppercase letter", "+length at least 16 characters"}},
		// Readme examples
		{"something", nil, 1, nil},
		{"somethin1", nil, 2, nil},
		{"somethingnew", nil, 2, nil},
		{"Somethin1", nil, 3, nil},
		{"somethinglonger", nil, 3, nil},
		{"Someth!n1", nil, 4, nil},
		{"somethingmuchlonger", nil, 4, nil},
		{"Someth!n10", nil, 5, nil},
		{"correct horse battery staple", nil, 5, nil},
	} {
		score, info := Check(test.password, test.extraWords)
		sort.Strings(info)
		if score != test.expectedScore {
			t.Errorf("%q: score = %v (%#v), want %v (%#v)", test.password, score, score, test.expectedScore, test.expectedScore)
		}
		if test.expectedInfo != nil && !reflect.DeepEqual(info, test.expectedInfo) {
			t.Errorf("%q: info = %#v, want %#v", test.password, info, test.expectedInfo)
		}
	}
}

func TestScoreNames(t *testing.T) {
	for _, test := range []struct {
		score Score
		name  string
	}{
		{-10, "Terrible"},
		{-1, "Terrible"},
		{1, "Weak"},
		{3, "Okay"},
		{4, "Good"},
		{10, "Strong"},
	} {
		if s := test.score.String(); s != test.name {
			t.Errorf("String = %v, want %v", s, test.name)
		}
	}
}

func TestExtractWords(t *testing.T) {
	for _, test := range []struct {
		word    string
		extract []string
	}{
		{"foo-bar", []string{"foo", "bar"}},
		{"!!", []string{}},
		{"foo@example.co", []string{"foo", "example"}},
		{"🤓🤓🤓🤓🤓🤓🤓", []string{}},
		{"日本", []string{"日本"}},
	} {
		e := ExtractWords(test.word)
		if !reflect.DeepEqual(test.extract, e) {
			t.Errorf("%q: e = %#v, want %#v", test.word, e, test.extract)
		}
	}
}

func BenchmarkCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Check("some-ok-thing2", []string{"foo"})
	}
}
