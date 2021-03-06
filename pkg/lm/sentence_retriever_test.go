package lm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/suggest-go/suggest/pkg/alphabet"
)

func TestSentenceRetrieve(t *testing.T) {
	cases := []struct {
		text     string
		expected []Sentence
	}{
		{
			"i wanna rock. hello my friend. what? dab. чтоооо. ты - не я",
			[]Sentence{
				{"i", "wanna", "rock"},
				{"hello", "my", "friend"},
				{"what"},
				{"dab"},
				{"чтоооо"},
				{"ты", "не", "я"},
			},
		},
	}

	tokenizer := NewTokenizer(alphabet.NewCompositeAlphabet(
		[]alphabet.Alphabet{
			alphabet.NewEnglishAlphabet(),
			alphabet.NewRussianAlphabet(),
			alphabet.NewNumberAlphabet(),
		},
	))

	stopAlphabet := alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'})

	for _, c := range cases {
		retriever := NewSentenceRetriever(
			tokenizer,
			strings.NewReader(c.text),
			stopAlphabet,
		)

		actual := []Sentence{}

		for s := retriever.Retrieve(); s != nil; s = retriever.Retrieve() {
			actual = append(actual, s)
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("Test fail, expected %v, got %v", c.expected, actual)
		}
	}
}
