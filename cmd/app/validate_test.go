package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"RupenderSinghRathore/AuthCli/internal/validator"
)

func TestValidateUsernamePassword(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  []validator.ValidationError
	}{
		{
			name:     "valid username and password",
			username: "validuser",
			password: "validpassword",
			wantErr:  nil,
		},
		{
			name:     "invalid username and password",
			username: "$$$%$%^",
			password: "12",
			wantErr: []validator.ValidationError{
				{Field: "username", Message: "must be 3-32 characters and contain only letters, numbers, '_' or '-'"},
				{Field: "password", Message: "must be at least 4 characters"},
			},
		},
		{
			name:     "empty username and password",
			username: "",
			password: "",
			wantErr: []validator.ValidationError{
				{Field: "username", Message: "required"},
				{Field: "password", Message: "must be at least 4 characters"},
			},
		},
		{
			name:     "too short username",
			username: "in",
			password: "validpassword",
			wantErr: []validator.ValidationError{
				{
					Field:   "username",
					Message: "must be 3-32 characters and contain only letters, numbers, '_' or '-'",
				},
			},
		},
		{
			name:     "invalid characters in username",
			username: "i++**",
			password: "validpassword",
			wantErr: []validator.ValidationError{
				{
					Field:   "username",
					Message: "must be 3-32 characters and contain only letters, numbers, '_' or '-'",
				},
			},
		},
		{
			name:     "too short password",
			username: "validuser",
			password: "123",
			wantErr: []validator.ValidationError{
				{Field: "password", Message: "must be at least 4 characters"},
			},
		},
		{
			name:     "too long password",
			username: "validuser",
			password: strings.Repeat("a", 73),
			wantErr: []validator.ValidationError{
				{Field: "password", Message: "must be at most 72 characters"},
			},
		},
		{
			name:     "boundary case: valid/lower-bound",
			username: "123",
			password: "1234",
			wantErr:  nil,
		},
		{
			name:     "boundary case: valid/upper-bound",
			username: strings.Repeat("a", 32),
			password: strings.Repeat("a", 72),
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUsernamePassword(tt.username, []byte(tt.password))
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}

			vErr, ok := errors.AsType[*validator.Validator](err)
			if !ok {
				t.Fatalf("expected *validator.Validator, got %T", err)
			}
			if !reflect.DeepEqual(vErr.Errors, tt.wantErr) {
				t.Fatalf(
					"got errors %#v, want %#v",
					vErr.Errors,
					tt.wantErr,
				)
			}
		})
	}
}
