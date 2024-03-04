package entities

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	type args struct {
		user     User
		lifetime time.Duration
		algo     string
		secret   string
	}

	type result struct {
		token     string
		expiredAt time.Time
		err       error
	}

	lifetime := time.Duration(2)

	tests := []struct {
		name   string
		args   args
		wanted result
	}{
		{
			name: "Invalid algo",
			args: args{
				user:     User{},
				lifetime: lifetime,
				algo:     "",
				secret:   "my-secret",
			},
			wanted: result{
				token:     "",
				expiredAt: time.Now(),
				err:       errors.New("unsupported JWT algo: must be HS512 or ES384"),
			},
		},
		{
			name: "Invalid algo",
			args: args{
				user:     User{},
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "secret",
			},
			wanted: result{
				token:     "",
				expiredAt: time.Now(),
				err:       errors.New("secret must have at least 8 characters"),
			},
		},
		{
			name: "Valid",
			args: args{
				user:     User{},
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "my-secret",
			},
			wanted: result{
				token:     "",
				expiredAt: time.Now().Add(lifetime * time.Hour),
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiredAt, err := tt.args.user.GenerateJWT(
				tt.args.lifetime,
				tt.args.algo,
				tt.args.secret,
			)
			got := result{token, expiredAt, err}

			if got.err != nil {
				assert.Equal(t, got.token, tt.wanted.token)
			} else {
				assert.Greater(t, len(got.token), 0)
				assert.Greater(t, got.expiredAt, time.Now().Add(lifetime*time.Hour-time.Minute))
				assert.Less(t, got.expiredAt, time.Now().Add(lifetime*time.Hour+time.Minute))
			}
			assert.Equal(t, got.err, tt.wanted.err)
		})
	}
}
