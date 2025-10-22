package words

import (
	"bufio"
	"embed"
	"math/rand"
	"time"
)

// Rng returns the underlying random number generator
func (w *Words) Rng() *rand.Rand {
	return w.rng
}

//go:embed data/dict.txt
var dictFile embed.FS

type Words struct {
	rng *rand.Rand
}

func New() *Words {
	src := rand.NewSource(time.Now().UnixNano())
	return &Words{rng: rand.New(src)}
}

// GetOne returns a random word from the embedded dictionary using reservoir sampling (streaming, memory efficient)
func (w *Words) GetOne() string {
	f, err := dictFile.Open("data/dict.txt")
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var result string
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		count++
		if w.rng.Intn(count) == 0 {
			result = line
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return result
}

// GetMany returns n unique random words (streaming, memory efficient)
func (w *Words) GetMany(n int) []string {
	if n <= 0 {
		return make([]string, 0)
	}

	result := make([]string, 0, n)
	seen := make(map[string]struct{})
	for len(result) < n {
		word := w.GetOne()
		if word == "" {
			continue
		}
		if _, exists := seen[word]; !exists {
			seen[word] = struct{}{}
			result = append(result, word)
		}
	}
	return result
}
