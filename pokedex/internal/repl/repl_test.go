package repl
import (
	"testing"
	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	} {
		{
			input: "   hello world   ",
			expected: []string{"hello", "world"},
		},
	}
	for _, useCase := range cases {
		result := cleanInput(useCase.input)
		if !cmp.Equal(result, useCase.expected) {
			t.Fatalf("expected: %#v got: %#v", useCase.expected, result)
		}
	}
}
