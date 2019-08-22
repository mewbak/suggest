package cmd

import (
	"bufio"
	"fmt"
	"github.com/alldroll/suggest/pkg/store"
	"os"

	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(countNGramsCmd)
}

var countNGramsCmd = &cobra.Command{
	Use:   "ngram-count -c [config path]",
	Short: "builds ngram counts for the given config file using google ngram format",
	Long:  `builds ngram counts for the given config file using google ngram format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lm.ReadConfig(configPath)

		if err != nil {
			return fmt.Errorf("could read config %s", err)
		}

		trie, err := buildNGramsCount(config)

		if err != nil {
			return err
		}

		return storeNGramsCount(config, trie)
	},
}

// buildNGramsCount builds a count trie
func buildNGramsCount(config *lm.Config) (lm.CountTrie, error) {
	sourceFile, err := os.Open(config.GetSourcePath())

	if err != nil {
		return nil, fmt.Errorf("could read source file %s", err)
	}

	defer sourceFile.Close()

	retriever := lm.NewSentenceRetriever(
		lm.NewTokenizer(config.GetWordsAlphabet()),
		bufio.NewReader(sourceFile),
		config.GetSeparatorsAlphabet(),
	)

	builder := lm.NewNGramBuilder(
		config.StartSymbol,
		config.EndSymbol,
	)

	return builder.Build(retriever, config.NGramOrder), nil
}

// storeNGramsCount flushes the constructed count trie on FS
func storeNGramsCount(config *lm.Config, trie lm.CountTrie) error {
	directory, err := store.NewFSDirectory(config.GetOutputPath())

	if err != nil {
		return fmt.Errorf("failed to create a fs directory: %v", err)
	}

	writer := lm.NewGoogleNGramWriter(config.NGramOrder, directory)

	if err := writer.Write(trie); err != nil {
		return fmt.Errorf("could save ngrams %s", err)
	}

	return nil
}
