package password

import (
	"errors"
	"fmt"
	"testing"
	"testing/iotest"

	"github.com/tecnickcom/gogen/pkg/random"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

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
