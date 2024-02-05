package values_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	type result struct {
		password Password
		err      error
	}
	tests := []struct {
		value  string
		wanted result
	}{
		{
			value: "password",
			wanted: result{
				password: Password{value: "password"},
				err:      nil,
			},
		},
		{
			value: "bad",
			wanted: result{
				password: Password{},
				// TODO: Fix this
				err: nil, // utils.ValidatorErrors[{"value": "min=8"}],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := NewPassword(tt.value)

			if err != nil {
				assert.Equal(t, err, tt.wanted.err)
			} else {
				assert.Equal(t, got, tt.wanted.password)
			}
		})
	}
}
