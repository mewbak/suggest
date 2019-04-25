package lm

import (
	"math"
	"os"
	"testing"
)

func TestScoreSentenceFromFile(t *testing.T) {
	config := Config{
		NGramOrder:  3,
		StartSymbol: "<S>",
		EndSymbol:   "</S>",
		OutputPath:  "testdata/fixtures",
	}

	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	reader := NewGoogleNGramReader(config.NGramOrder, indexer, config.OutputPath)

	model, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	lm := NewLanguageModel(model, indexer, &config)
	testLM(lm, t)
}

func TestScoreSentenceFromBinary(t *testing.T) {
	f, err := os.Open("testdata/config-example.json")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	config, err := ReadConfig(f)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	lm, err := RetrieveLMFromBinary(config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	testLM(lm, t)
}

func testLM(lm LanguageModel, t *testing.T) {
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
		actual, _ := lm.ScoreSentence(c.sentence)

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
