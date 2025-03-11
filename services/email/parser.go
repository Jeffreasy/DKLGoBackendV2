package email

import (
	"strings"
)

// ParseEmailContent extracts text content from email body
func (s *EmailService) ParseEmailContent(content string) string {
	// Process HTML content first
	content = s.ProcessHTML(content)

	// Clean up whitespace
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\n\n\n", "\n\n")

	return content
}
