package utils

import (
	"os"
	"testing"
	"wordcrunch/internal/models"
)

func TestCountWordStarts(t *testing.T) {
	cases := []struct {
		name        string
		buf         string
		prevIsSpace bool
		want        int
	}{
		{"empty", "", true, 0},
		{"single word", "hello", true, 1},
		{"leading space", " hello", true, 1},
		{"only spaces", "   \t\n", true, 0},
		{"multiple words", "hello world foo", true, 3},
		{"prev not space, continues word", "llo world", false, 1},
		{"prev not space, starts with space", " world", false, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := countWordStarts([]byte(c.buf), c.prevIsSpace)
			if got != c.want {
				t.Errorf("got %d, want %d", got, c.want)
			}
		})
	}
}

func TestProcessChunk(t *testing.T) {
	content := "hello world\nfoo bar baz"
	f, err := os.CreateTemp("", "wordcrunch-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}

	// Split into two chunks: [0,12) and [12, len)
	c1 := models.Chunk{Offset: 0, Length: 12}
	c2 := models.Chunk{Offset: 12, Length: int64(len(content)) - 12}

	n1, err := ProcessChunk(f, c1)
	if err != nil {
		t.Fatal(err)
	}
	n2, err := ProcessChunk(f, c2)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := n1+n2, 5; got != want {
		t.Errorf("total words = %d, want %d", got, want)
	}
}
