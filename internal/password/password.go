// Package password generates random passwords.
package password

import (
	"fmt"

	"github.com/tecnickcom/gogen/pkg/random"
)

// Password contains the random generator configuration.
type Password struct {
	Charset  string `json:"charset"  validate:"required,min=1,max=256,rndcharset"`
	Length   int    `json:"length"   validate:"required,min=1,max=4096"`
	Quantity int    `json:"quantity" validate:"required,min=1,max=1000"`
	rnd      *random.Rnd
}

// New instantiate a new Password generator object.
func New(charset string, length, quantity int) *Password {
	// Duplicate characters would bias the output toward them, so the effective
	// charset only keeps the first occurrence of each character.
	charset = dedupCharset(charset)

	return &Password{
		Charset:  charset,
		Length:   length,
		Quantity: quantity,
		rnd:      random.New(nil, random.WithByteToCharMap([]byte(charset))),
	}
}

// dedupCharset removes duplicate bytes from the charset while preserving the
// order of first appearance.
func dedupCharset(charset string) string {
	var seen [256]bool

	out := make([]byte, 0, len(charset))

	for i := range len(charset) {
		c := charset[i]
		if seen[c] {
			continue
		}

		seen[c] = true

		out = append(out, c)
	}

	return string(out)
}

// Generate returns the specified amount of random passwords.
func (p *Password) Generate() ([]string, error) {
	lst := make([]string, p.Quantity)

	for i := range p.Quantity {
		s, err := p.rnd.RandString(p.Length)
		if err != nil {
			return nil, fmt.Errorf("failed generating random password: %w", err)
		}

		lst[i] = s
	}

	return lst, nil
}
