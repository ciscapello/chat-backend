package userhandler

import "testing"

func TestIsValidEmail(t *testing.T) {
	uh := UserHandler{}

	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"invalid.email", false},
		{"invalidemail", false},
		{"invali@demail", false},
		{"user@domain", false},
		{"user@domain.c", false},
		{"user@domain.com", true},
		{"user.name@domain.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := uh.isValidEmail(tt.email); got != tt.want {
				t.Errorf("UserHandler.isValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
