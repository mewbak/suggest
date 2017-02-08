package suggest

const (
	LEVENSHTEIN = iota
	NGRAM
	JACCARD
)

type EditDistance interface {
	Calc(a, b string) float64
}

func GetEditDistance(t int, k int) (EditDistance, error) {
	switch t {
	case LEVENSHTEIN:
		return &LevenshteinDistance{}, nil

	case NGRAM:
		return &NGramDistance{k}, nil

	case JACCARD:
		return &JaccardDistance{k}, nil

	default:
		//return nil, errors.New("Invalid type")
		return nil, nil
	}
}

type LevenshteinDistance struct{}

func (self *LevenshteinDistance) Calc(a, b string) float64 {
	aLen, bLen := len(a), len(b)
	if aLen == 0 {
		return float64(bLen)
	}

	if bLen == 0 {
		return float64(aLen)
	}

	r1, r2 := []rune(a), []rune(b)
	column := make([]int, aLen+1)

	for i := 1; i < aLen+1; i++ {
		column[i] = i
	}

	for j := 1; j < bLen+1; j++ {
		column[0] = j
		prev := j - 1
		for i := 1; i < aLen+1; i++ {
			tmp := column[i]
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			column[i] = min3(
				column[i]+1,
				column[i-1]+1,
				prev+cost,
			)
			prev = tmp
		}
	}

	return float64(column[aLen])
}

type NGramDistance struct {
	k int
}

/*
 * ngram distance between a, b defined in
 * "Approximate string-matching with q-grams and maximal matches"
 *
 * Complexity O(aLen + bLen)
 */
func (self *NGramDistance) Calc(a, b string) float64 {
	profileA, profileB := getProfile(a, self.k), getProfile(b, self.k)
	union := make(map[string]byte)
	for k := range profileA {
		union[k] = 1
	}

	for k := range profileB {
		union[k] = 1
	}

	distance := 0.0
	for key, _ := range union {
		freqA, freqB := 0, 0
		if val, ok := profileA[key]; ok {
			freqA = val
		}

		if val, ok := profileB[key]; ok {
			freqB = val
		}

		d := freqA - freqB
		if d < 0 {
			d = -d
		}

		distance += float64(d)
	}

	return distance
}

type JaccardDistance struct {
	k int
}

func (self *JaccardDistance) Calc(a, b string) float64 {
	if a == b {
		return 0
	}

	profileA, profileB := getProfile(a, self.k), getProfile(b, self.k)
	union := make(map[string]byte)
	for k := range profileA {
		union[k] = 1
	}

	for k := range profileB {
		union[k] = 1
	}

	inter := 0
	for key, _ := range union {
		part := 1
		if _, ok := profileA[key]; !ok {
			part = 0
		}

		if _, ok := profileB[key]; !ok {
			part = 0
		}

		inter += part
	}

	return 1.0 - float64(inter)/float64(len(union))
}