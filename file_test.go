package learning

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	file, err := os.Open("name.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

func TestReadFile(t *testing.T) {
	byteFile, err := os.ReadFile("name.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(byteFile))
}

func TestMin(t *testing.T) {
	a := min(1, 2)
	b := max(1, 2)
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
}
