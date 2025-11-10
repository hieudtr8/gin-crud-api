package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		// Valid emails
		{
			name:  "Simple valid email",
			email: "user@example.com",
			want:  true,
		},
		{
			name:  "Email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "Email with numbers",
			email: "user123@example456.com",
			want:  true,
		},
		{
			name:  "Email with dots in local part",
			email: "first.last@example.com",
			want:  true,
		},
		{
			name:  "Email with plus sign",
			email: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "Email with hyphen in domain",
			email: "user@my-domain.com",
			want:  true,
		},
		{
			name:  "Email with underscore",
			email: "user_name@example.com",
			want:  true,
		},
		{
			name:  "Email with percent",
			email: "user%tag@example.com",
			want:  true,
		},
		{
			name:  "Long TLD",
			email: "user@example.technology",
			want:  true,
		},
		{
			name:  "Two-letter TLD",
			email: "user@example.io",
			want:  true,
		},

		// Invalid emails
		{
			name:  "Missing @ symbol",
			email: "userexample.com",
			want:  false,
		},
		{
			name:  "Missing local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "Missing domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "Missing TLD",
			email: "user@example",
			want:  false,
		},
		{
			name:  "Empty string",
			email: "",
			want:  false,
		},
		{
			name:  "Multiple @ symbols",
			email: "user@@example.com",
			want:  false,
		},
		{
			name:  "Spaces in email",
			email: "user @example.com",
			want:  false,
		},
		{
			name:  "TLD too short (1 char)",
			email: "user@example.c",
			want:  false,
		},
		{
			name:  "Special chars in domain",
			email: "user@exam ple.com",
			want:  false,
		},
		{
			name:  "Missing domain name",
			email: "user@.com",
			want:  false,
		},
		// Note: Edge cases like double dots, leading/trailing dots are not caught by the simple regex
		// This is acceptable for basic validation - more complex validation would require RFC 5322 parser
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEmail(tt.email)
			assert.Equal(t, tt.want, got, "isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
		})
	}
}
