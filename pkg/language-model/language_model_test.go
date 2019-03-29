package lm

import (
	"math"
	"testing"
)

func TestScoreSentenceFromFile(t *testing.T) {
	indexer := NewIndexer()
	reader := NewGoogleNGramReader(3, indexer, "fixtures")

	model, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	generator := NewGenerator(3, start, end)
	lm := NewLanguageModel(model, generator, indexer)

	cases := []struct {
		sentence      Sentence
		expectedScore float64
	}{
		{Sentence{"i", "am", "sam"}, -1.3862},
		{Sentence{"i", "am"}, -1.3862},
		{Sentence{"sam", "i", "am"}, -0.6931},
		{Sentence{"sam", "am", "i"}, -10.2852},
		{Sentence{"i", "dont", "know"}, -105.0514},
		{Sentence{"no", "one", "word"}, -203.7297},
	}

	for _, c := range cases {
		actual := lm.ScoreSentence(c.sentence)

		if diff := math.Abs(actual - c.expectedScore); diff >= tolerance {
			t.Errorf(
				"Test fail, for %v expected score %v, got %v",
				c.sentence,
				c.expectedScore,
				actual,
			)
		}
	}
}
