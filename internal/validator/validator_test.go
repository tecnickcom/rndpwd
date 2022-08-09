package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	t.Parallel()

	type valTest struct {
		Charset string `json:"charset" validate:"rndcharset"`
	}

	tests := []struct {
		name    string
		in      *valTest
		wantErr bool
	}{
		{
			name:    "valid",
			in:      &valTest{Charset: "valid"},
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      &valTest{Charset: "in va lid"},
			wantErr: true,
		},
		{
			name:    "empty",
			in:      &valTest{Charset: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v, err := New("json")
			require.NoError(t, err)

			err = v.ValidateStruct(tt.in)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
