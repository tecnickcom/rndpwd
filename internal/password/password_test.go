package password

import (
	"fmt"
	"testing"

	"github.com/tecnickcom/rndpwd/internal/validator"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	length := 64
	quantity := 3

	p := New(validator.ValidCharset, length, quantity)

	pwds := p.Generate()

	if q := len(pwds); q != quantity {
		t.Error(fmt.Errorf("The expected quantity is %d, found %d", quantity, q))
	}

	for _, pwd := range pwds {
		if l := len(pwd); l != length {
			t.Error(fmt.Errorf("The expected password length is %d, found %d", length, l))
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	p := New(validator.ValidCharset, 32, 1)

	b.ResetTimer()

	for range b.N {
		p.Generate()
	}
}
