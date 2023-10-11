// Package password generates random passwords.
package password

import (
	"crypto/rand"
	"math/big"
)

// Password contains the random generator configuration.
type Password struct {
	Charset  string `json:"charset"  validate:"required,min=2,max=92,rndcharset"`
	Length   int    `json:"length"   validate:"required,min=2"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

// New instantiate a new Password generator object.
func New(charset string, length, quantity int) *Password {
	return &Password{
		Charset:  charset,
		Length:   length,
		Quantity: quantity,
	}
}

// Generate returns the specified amount of random passwords.
func (p *Password) Generate() []string {
	lst := make([]string, p.Quantity)

	for i := 0; i < p.Quantity; i++ {
		lst[i] = p.newPassword()
	}

	return lst
}

// newPassword returns a random password containing "length" characters from the "charset" string.
func (p *Password) newPassword() string {
	pwd := make([]byte, p.Length)
	chars := []byte(p.Charset)
	maxValue := new(big.Int).SetInt64(int64(len(p.Charset)))

	for i := 0; i < p.Length; i++ {
		rnd, err := rand.Int(rand.Reader, maxValue)
		if err == nil {
			pwd[i] = chars[rnd.Int64()]
		}
	}

	return string(pwd)
}
