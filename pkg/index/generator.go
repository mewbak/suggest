package index

const maxN = 8

// Generator represents entity for creating terms from given word
type Generator interface {
	// Generate returns terms array for given word
	Generate(word string) []Term
}

// generatorImpl implements Generator interface
type generatorImpl struct {
	nGramSize int
}

// NewGenerator returns new instance of Generator
func NewGenerator(nGramSize int) Generator {
	return &generatorImpl{
		nGramSize: nGramSize,
	}
}

// Generate returns terms array for given word
func (g *generatorImpl) Generate(word string) []Term {
	nGrams := SplitIntoNGrams(word, g.nGramSize)
	l := len(nGrams)
	set := make(map[Term]struct{}, l)
	list := make([]Term, 0, l)

	for _, nGram := range nGrams {
		if _, found := set[nGram]; !found {
			set[nGram] = struct{}{}
			list = append(list, nGram)
		}
	}

	return list
}

// SplitIntoNGrams split given word on k-gram list
func SplitIntoNGrams(word string, k int) []string {
	runes := []rune(word)
	result := make([]string, 0, len(runes))

	for i := 0; i < len(runes)-k+1; i++ {
		result = append(result, string(runes[i:i+k]))
	}

	return result
}
