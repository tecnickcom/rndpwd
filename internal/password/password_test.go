package password

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/tecnickcom/gogen/pkg/random"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

func TestDedupCharset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no duplicates", "abc", "abc"},
		{"adjacent duplicates", "aabbcc", "abc"},
		{"interleaved duplicates", "abcabc", "abc"},
		{"all same", "aaaa", "a"},
		{"preserves first-seen order", "cbacba", "cba"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := dedupCharset(tt.in); got != tt.want {
				t.Errorf("dedupCharset(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewDedupesCharset(t *testing.T) {
	t.Parallel()

	p := New("aaabbbccc", 16, 2)
	if p.Charset != "abc" {
		t.Errorf("expected deduped charset %q, found %q", "abc", p.Charset)
	}

	pwds, err := p.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, pwd := range pwds {
		for _, c := range pwd {
			if !strings.ContainsRune("abc", c) {
				t.Errorf("password %q contains %q outside the deduped charset", pwd, c)
			}
		}
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	length := 64
	quantity := 3

	p := New(validator.ValidCharset, length, quantity)

	pwds, err := p.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q := len(pwds); q != quantity {
		t.Error(fmt.Errorf("The expected quantity is %d, found %d", quantity, q))
	}

	for _, pwd := range pwds {
		if l := len(pwd); l != length {
			t.Error(fmt.Errorf("The expected password length is %d, found %d", length, l))
		}
	}
}

func TestGenerateError(t *testing.T) {
	t.Parallel()

	p := New(validator.ValidCharset, 16, 2)
	// swap the random source for one whose reader always fails
	p.rnd = random.New(iotest.ErrReader(errors.New("rng failure")), random.WithByteToCharMap([]byte(validator.ValidCharset)))

	pwds, err := p.Generate()
	if err == nil {
		t.Error("expected an error when the random reader fails")
	}

	if pwds != nil {
		t.Errorf("expected nil passwords on error, found %v", pwds)
	}
}

func BenchmarkGenerate(b *testing.B) {
	p := New(validator.ValidCharset, 32, 1)

	for b.Loop() {
		_, _ = p.Generate()
	}
}
