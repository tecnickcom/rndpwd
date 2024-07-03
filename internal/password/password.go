// Package password generates random passwords.
package password

import (
	"github.com/Vonage/gosrvlib/pkg/random"
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
	return &Password{
		Charset:  charset,
		Length:   length,
		Quantity: quantity,
		rnd:      random.New(nil, random.WithByteToCharMap([]byte(charset))),
	}
}

// Generate returns the specified amount of random passwords.
func (p *Password) Generate() []string {
	lst := make([]string, p.Quantity)

	for i := range p.Quantity {
		lst[i], _ = p.rnd.RandString(p.Length)
	}

	return lst
}
