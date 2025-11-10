package graph

import "regexp"

// ============================================================================
// Validation Helper Functions
// ============================================================================

// isValidEmail validates email format using a regular expression
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
