package email

import (
	"strings"
)

func (s *EmailService) ProcessContactForm(content string) (map[string]string, error) {
	// Process HTML content first
	content = s.ProcessHTML(content)

	// Extract form fields
	fields := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try different separators
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(line, "=", 2)
		}
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			fields[key] = value
		}
	}

	return fields, nil
}

func (s *EmailService) ExtractFormData(content string) (map[string]string, error) {
	// Clean HTML first
	content = s.ProcessHTML(content)

	// Extract form fields
	fields := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try different separators
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(line, "=", 2)
		}
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			fields[key] = value
		}
	}

	return fields, nil
}
